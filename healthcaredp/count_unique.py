from pprint import pprint

with open('data/healthcare_cleaned.csv') as f:
    lines = f.readlines()
    header = lines[0]
    lines = lines[1:]
    lines_split = [[*i.split(',')] for i in lines]
    names = [i[0] for i in lines_split]
    conditions = [i[4] for i in lines_split]
    test_results = [i[14] for i in lines_split]

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

test_results_maps = [(k, len(v)) for k, v in test_results_maps.items()]
test_results_maps.sort(key=lambda x: x[1], reverse=True)

print(f"Count Parameters: CountConditionsDp")
# Maximum number of times a privacy key appears in the dataset, therefore the max number of times a key can contribute to a single partition
# This value is the same for all count operations on this dataset
print(f"\tMaxValue = {names_count_tuples[0][1]}")
# Maximum number of conditions that a single privacy key contributes to
print(f"\tMaxPartitionsContributed = {conditions_count_tuples[0][1]}")

print(f"Count Parameters: CountTestResults")
print(f"\tMaxValue = {names_count_tuples[0][1]}")
# Maximum number of test results that a single privacy key contributes to
print(f"\tMaxPartitionsContributed = {test_results_maps[0][1]}")







