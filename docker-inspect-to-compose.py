import json
import sys
import argparse
import yaml

def convert_to_compose(container_info):
    # Extract relevant information
    container_name = container_info['Name'][1:]  # remove the leading /
    image_name = container_info['Config']['Image']
    labels = container_info['Config']['Labels'] or {}
    environment = container_info['Config']['Env'] or []
    ports = container_info['HostConfig']['PortBindings'] or {}
    volumes = container_info['Mounts'] or []
    
    # Prepare ports
    port_mappings = []
    for host_port, bindings in ports.items():
        for binding in bindings:
            port_mappings.append(f"{binding['HostPort']}:{host_port}")

    # Prepare volumes
    volume_mappings = []
    for mount in volumes:
        volume_mappings.append(f"{mount['Source']}:{mount['Destination']}")

    # Construct the Docker Compose structure
    compose = {
        'version': '3',
        'services': {
            container_name: {
                'image': image_name,
                'container_name': container_name,
                'ports': port_mappings,
                'environment': environment,
                'volumes': volume_mappings,
                'labels': labels,
            }
        }
    }

    return compose

def main():
    parser = argparse.ArgumentParser(description="Convert docker inspect output to docker-compose.yml")
    parser.add_argument('input', help="Path to JSON file or '-' for stdin")
    args = parser.parse_args()

    # Read the input data
    if args.input == '-':
        container_info = json.load(sys.stdin)
    else:
        with open(args.input, 'r') as f:
            container_info = json.load(f)

    # If multiple containers are inspected, take the first one
    if isinstance(container_info, list):
        container_info = container_info[0]

    # Convert and print to YAML
    compose = convert_to_compose(container_info)
    print(yaml.dump(compose, default_flow_style=False))

if __name__ == '__main__':
    main()