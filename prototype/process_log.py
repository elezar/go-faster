import re
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
    du = speed_matcher.search(text)
    dl, dlu = split_units(du.group(1))
    ul, ulu = split_units(du.group(2))
    err = 0
  except:
    dl = 0
    ul = 0
    err = 1

  # print dl, dlu, ul, ulu

  return dl, ul, err


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
    dl, ul, err = process_item(output)

    times.append(start_time)
    rates.append([dl, ul, err])

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

def process_log(filename):

  t, v = get_data(filename)
  valid = np.where(v[:,2] == 0)
  invalid = np.where(v[:,2] == 1)

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
  P.stem(tiv, 9*np.ones(tiv.size), linefmt='k-', markerfmt='ko', basefmt='k', label='Error (NO CONNECTION)')


  P.plot_date(t01, maxes[0]*y01, 'b:')
  P.plot_date(t01, maxes[1]*y01, 'r:')

  P.yticks(np.arange(0,55,5))
  P.ylim([0, 51])
  P.ylabel('Rate [Mbit/s]')
  P.legend(loc='best')

  return f

if __name__ == "__main__":
  import sys
  try:
    logname = sys.argv[1]
  except:
    pass

  process_log(logname)
  P.show()
