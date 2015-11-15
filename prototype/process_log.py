import re
import argparse
import requests
import json
import numpy as np
import numpy.ma as ma
import pylab as P
import datetime

item_matcher = re.compile("^([ \d:/]+) ### Starting:(.*?)^([ \d:/]+) ### Done:$", re.S+re.M)
source_matcher = re.compile("Testing from (.*?)\s+\(([\d\.]+)\)", re.S)
dest_matcher = re.compile("Hosted by (.*?)\[(.*?)\]:\s+(.*?)$", re.M+re.S)
speed_matcher = re.compile("^Download: ([\d\.]+\s+.*?)$.*Upload: ([\d\.]+\s+.*?)$", re.M+re.S)

recognised_units = ['s', 'm', 'bit/s', ]

prefix_map = {'m': 1.e-3, 'k': 1.0e3, 'M': 1.0e6, 'G': 1.0e9 }

dl_col  = 0
ul_col = 1
distance_col = 2
latency_col = 3
err_col = 4

def get_factor(prefix):
  if prefix in prefix_map:
    return prefix_map[prefix], ''

  return 1.0, 'prefix'

def normalise(units):
  if len(units) == 1 or units in recognised_units:
    return 1.0, units

  factor, prefix = get_factor(units[0])

  return factor, prefix+units[1:]


def split_units(string):
  value, units = string.split(' ')

  factor, units = normalise(units)

  return factor*float(value), units


def process_item(text):


  try:
    s = source_matcher.search(text)
    d = dest_matcher.search(text[s.end():])
    host = d.group(1)
    distance_m, distance_units = split_units(d.group(2))
    latency_s, latency_units = split_units(d.group(3))
    du = speed_matcher.search(text)
    dl_bs, dlu = split_units(du.group(1))
    ul_bs, ulu = split_units(du.group(2))
    err = 0
  except:
    distance_m = 0
    latency_s = 0
    dl_bs = 0
    ul_bs = 0
    err = 1

  # print dl, dlu, ul, ulu

  return dl_bs, ul_bs, distance_m, latency_s, err


def set_axis():
  print


def get_data(filename):
  f = open(filename, "r")
  text = f.read()
  f.close()

  times = []
  rates = []

  for m in item_matcher.finditer(text):
    start_time = datetime.datetime.strptime(m.group(1), '%Y/%m/%d %H:%M:%S')
    output = m.group(2)
    end_time = m.group(3)
    dl, ul, distance_m, latency_s, err = process_item(output)

    times.append(start_time)
    rates.append([dl, ul, distance_m, latency_s, err])

  v = np.array(rates)

  t = np.array(times)

  return t, v

def get_average(t, v):
  ts = np.empty(t.shape)

  for i in range(len(t)):
    dt = t[i] - t[0]
    ts[i] = dt.days*24*60*60.0 + dt.seconds

  average = np.trapz(v, ts)/(ts[-1])

  return average



def populate_api(filename, post_url):
  """Read the specified API and populate the database from the logfile."""

  t, v = get_data(filename)

  num_data = len(t)

  import time
  headers = {'content-type': 'application/json'}

  # proxies = None
  proxies = {
    "http": "http://go-faster.devfest.com:2000",
  }
  for i in range(num_data):
    data_point = {}
    dt = t[i]
    data_point['unixtime'] = int(time.mktime(dt.timetuple()))
    data_point['download_speed_bs'] = int(v[i,dl_col])
    data_point['upload_speed_bs'] = int(v[i,ul_col])
    data_point['distance_m'] = int(v[i,distance_col])
    data_point['latency_ms'] = int(v[i,latency_col]*1000)

    data = json.dumps(data_point)

    # print data
    response = requests.put(post_url, data=data, proxies=proxies, headers=headers)

    if response.status_code != 200:
      print response



def process_log(filename):

  t, v = get_data(filename)
  valid = np.where(v[:,err_col] == 0)
  invalid = np.where(v[:,err_col] == 1)

  tv = t[valid]
  tiv = t[invalid]


  # P.plot_date(tiv, 4.5*np.ones(tiv.size), 'ko')
  # P.bar(tiv, 4.5*np.ones(tiv.size), 'ko')


  t01 = [t[0], t[-1]]
  y01 = np.ones(2)

  means = np.zeros(v.shape[1])
  for i in range(v.shape[1]):
    means[i] = get_average(t, v[:,i])/1e6

  # means = np.average(v, axis=0)/1e6
  maxes = np.array([50, 5])

  f = P.figure()
  P.plot_date(t, v[:,0]/1e6, 'b-', label='Download')
  P.plot_date(t01, means[0]*y01, 'b:', label='Average (%.0f Mbit/s : %.0f%%)' % (means[0], means[0]/maxes[0]*100))
  P.plot_date(t, v[:,1]/1e6, 'r-', label='Upload')
  P.plot_date(t01, means[1]*y01, 'r:', label='Average (%.0f Mbit/s : %.0f%%)' % (means[1], means[1]/maxes[1]*100))

  if len(tiv) > 0:
    P.stem(tiv, 9*np.ones(tiv.size), linefmt='k-', markerfmt='ko', basefmt='k', label='Error (NO CONNECTION)')


  P.plot_date(t01, maxes[0]*y01, 'b:')
  P.plot_date(t01, maxes[1]*y01, 'r:')

  P.yticks(np.arange(0,55,5))
  P.ylim([0, 51])
  P.ylabel('Rate [Mbit/s]')
  P.legend(loc='best')

  return f

if __name__ == "__main__":

  parser = argparse.ArgumentParser(description='Process a log file')

  parser.add_argument('logname', nargs=1,
                     help='The name of the logfile to process')
  parser.add_argument('--show-plot', action='store_true', default=False,
                     help='A flag to indicate whether the DB should be updated from the log file')
  parser.add_argument('--series-name', default='demo')
  parser.add_argument('--root-url', default="http://go-faster.devfest.com:8080",
                     help='The name of the logfile to process')
  parser.add_argument('--update-db', action='store_true', default=False,
                     help='A flag to indicate whether the DB should be updated from the log file')

  args = parser.parse_args()

  logfile = args.logname[0]

  if args.update_db:
    api_url = "%s/series/%s/data" % (args.root_url, args.series_name)

    populate_api(logfile, api_url)
  else:
    process_log(logfile)
    if args.show_plot:
      P.show()
