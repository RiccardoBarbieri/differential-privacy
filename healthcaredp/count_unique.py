from pprint import pprint
from datetime import datetime
import csv
import pandas as pd

with open('data/healthcare_cleaned.csv') as f:
    lines = f.readlines()
    header = lines[0]
    lines = lines[1:]
    csv_reader = csv.reader(lines)
    lines_split = [i for i in csv_reader]
    names = [i[0] for i in lines_split]
    conditions = [i[4] for i in lines_split]
    test_results = [i[14] for i in lines_split]

    date_fmt = "%Y-%m-%dT%H:%M:%SZ"
    admissions = [datetime.strptime(i[5], date_fmt) for i in lines_split]
    discharges = [datetime.strptime(i[12], date_fmt) for i in lines_split]

# Privacy Key -> Name

# Counts the number of times that the privacy key appears in the dataset
names_map = {}
for name in names:
    if name in names_map:
        names_map[name] += 1
    else:
        names_map[name] = 1

# Counts the number of different conditions that a privacy key contributes to
conditions_map = {}
for name, condition in zip(names, conditions):
    if name in conditions_map:
        conditions_map[name].add(condition)
    else:
        conditions_map[name] = {condition}

test_results_maps = {}
for name, test_result in zip(names, test_results):
    if name in test_results_maps:
        test_results_maps[name].add(test_result)
    else:
        test_results_maps[name] = {test_result}


names_count_tuples = [(k, v) for k, v in names_map.items()]
names_count_tuples.sort(key=lambda x: x[1], reverse=True)

conditions_count_tuples = [(k, len(v)) for k, v in conditions_map.items()]
conditions_count_tuples.sort(key=lambda x: x[1], reverse=True)

test_results_tuples = [(k, len(v)) for k, v in test_results_maps.items()]
test_results_tuples.sort(key=lambda x: x[1], reverse=True)

print(f"Count Parameters: CountConditionsDp")
# Maximum number of times a privacy key appears in the dataset, therefore the max number of times a key can contribute to a single partition
# This value is the same for all count operations on this dataset
print(f"\tMaxValue = {names_count_tuples[0][1]}")
# Maximum number of conditions that a single privacy key contributes to
print(f"\tMaxPartitionsContributed = {conditions_count_tuples[0][1]}")

print(f"Count Parameters: CountTestResults")
print(f"\tMaxValue = {names_count_tuples[0][1]}")
# Maximum number of test results that a single privacy key contributes to
print(f"\tMaxPartitionsContributed = {test_results_tuples[0][1]}")

durations = [dis - adm for adm, dis in zip(admissions, discharges)]
week_duration_tuples = [(adm.isocalendar()[1], duration) for duration, adm in zip(durations, admissions)]
df = pd.DataFrame(week_duration_tuples, columns=['week of year', 'duration'])
groups_mean = df.groupby(by="week of year").mean()
# pprint(groups_mean)

admission_week_map = {}
week_admission_map = {}
for name, admission in zip(names, admissions):
    week_of_year = admission.isocalendar()[1]
    if name in admission_week_map:
        admission_week_map[name].add(week_of_year)
    else:
        admission_week_map[name] = {week_of_year}

    if week_of_year in week_admission_map:
        week_admission_map[week_of_year].add(name)
    else:
        week_admission_map[week_of_year] = {name}

admission_week_tuples = [(k, len(v)) for k, v in admission_week_map.items()]
admission_week_tuples.sort(key=lambda x: x[1], reverse=True)

week_admission_tuples = [(k, v) for k, v in week_admission_map.items()]
week_admission_tuples.sort(key=lambda x: x[1], reverse=True)



print(f"Mean Parameters: MeanStayByWeek")
print(f"\tMaxValue = {names_count_tuples[0][1]}")
# Maximum number of week of year that a single privacy key contributes to
print(f"\tMaxPartitionsContributed = {admission_week_tuples[0][1]}")
# Maximum number of privacy keys that contribute to a single week of year
print(f"\tMaxContributionsPerPartition = {week_admission_tuples[0][1]}")








