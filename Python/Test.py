#!/usr/bin/python

import time

RulesFile = 'Ribbon.xml'

import LSystem

if __name__ == "__main__":
    tree = open(RulesFile).read()
    #shapes = LSystem.Evaluate(tree, seed = 29)
    lsystem = LSystem.LSystem(tree)
    startTime = time.clock()
    shapes = lsystem.evaluate(tree)
    print time.clock() - startTime, " seconds"
