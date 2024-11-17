import argparse
import csv
import datetime
import math

from matplotlib import ticker
from matplotlib.pyplot import figure
import matplotlib.pyplot as plt
from matplotlib.ticker import AutoMinorLocator


def main():
    argparser = argparse.ArgumentParser()
    argparser.set_defaults(cmd=lambda: None, cmd_args=lambda x: [])
    argparser.add_argument('path', type=str, metavar="PATH")
    argparser.add_argument('output', type=str, metavar="OUTPUT")

    args = argparser.parse_args()

    csvfile = open(args.path, 'r', newline='')
    wr = csv.reader(csvfile,delimiter=',')

    capacity_series = []
    voltage_series = []

    min_voltage = 500
    max_voltage = 0
    min_capacity = 0

    entries_i = iter(wr)
    next(entries_i)
    for entry in entries_i:
        voltage = float(entry[1])
        capacity = float(entry[3]) * 1000
        voltage_series.append(voltage)
        capacity_series.append(capacity)

        min_voltage = math.floor(min(min_voltage, voltage))
        max_voltage = math.ceil(max(max_voltage, voltage))

    max_capacity = capacity_series[-1]
    v_scale_mult = 1
    fig, (ax1_voltage) = plt.subplots(1, 1)
    fig.set_size_inches(8.5, 4)

    ax1_voltage.grid(True)
    ax1_voltage.set_xlabel("Capacity [mAh]", color='black')

    # Voltage
    ax1_voltage.plot(capacity_series, voltage_series, color='green')
    ax1_voltage.yaxis.set_major_formatter(ticker.FormatStrFormatter('%g V'))
    ax1_voltage.set_ylim([min_voltage * v_scale_mult, max_voltage * v_scale_mult])
    ax1_voltage.set_ylabel("Voltage", color='green')

    plt.savefig(args.output)
main()