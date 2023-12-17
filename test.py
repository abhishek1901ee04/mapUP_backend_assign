import random
import json

def generate_random_list():
    numbers = list(range(1, 1000000))
    random.shuffle(numbers)
    return numbers

def save_to_json(numbers, output_file='random_numbers.json'):
    data = {'numbers': numbers}
    with open(output_file, 'w') as json_file:
        json.dump(data, json_file, indent=4)

if __name__ == '__main__':
    random_numbers = generate_random_list()
    save_to_json(random_numbers)
    print("Generated random list and saved to 'random_numbers.json'")