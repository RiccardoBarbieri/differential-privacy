import os
import random
import csv


def generate_random_data():
    data = []
    for i in range(10):
        data.append((i, random.randint(0, 1)))
    return data


data = generate_random_data()

read_data = []
if not os.path.exists('data.csv'):
    with open('data.csv', 'w', newline='') as f:
        writer = csv.DictWriter(f, delimiter=',', fieldnames=['subject', 'property'])
        writer.writeheader()
        for item in data:
            writer.writerow({'subject': item[0], 'property': item[1]})
else:
    with open('data.csv', 'r', newline='') as f:
        reader = csv.DictReader(f)
        for row in reader:
            read_data.append(row)



