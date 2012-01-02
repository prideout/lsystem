#!/usr/bin/python

RulesFile = 'Ribbon.xml'

import LSystem

if __name__ == "__main__":
    tree = open(RulesFile).read()
    #shapes = LSystem.Evaluate(tree, seed = 29)
    lsystem = LSystem.LSystem(tree)
    shapes = lsystem.evaluate(tree)
