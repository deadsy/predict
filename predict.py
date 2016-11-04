#!/usr/bin/python

import os
import random

dataset_dir = 'dataset'
causal = False
random.seed((None, 12345678)[causal])

def split_datasets():
  """generate a training set and a testing set"""
  names = [f.split('.')[0] for f in os.listdir(dataset_dir) if f.endswith('html.gz')]
  random.shuffle(names)
  # split 2/3 training, 1/3 testing
  n = len(names) * 2 / 3
  return names[:n], names[n:]

def main():
  (training, testing) = split_datasets()
  print('%d files in the training set' % len(training))
  print('%d files in the testing set' % len(testing))

main()
