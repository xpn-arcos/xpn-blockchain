import sys
import json
import hashlib


def write_json_file (file_name, content):
    json_str = json.dumps(json.loads(content), indent=4)
    
    with open(file_name, "w") as file:
        file.write(json_str)



def get_file_hash (file_name):
    hash_obj = hashlib.new('sha256')

    with open(file_name, 'rb') as f:
        for block in iter(lambda: f.read(4096), b''):
            hash_obj.update(block)

    return hash_obj.hexdigest()


#Main
write_json_file (sys.argv[1], sys.argv[2])
print(get_file_hash (sys.argv[1]))