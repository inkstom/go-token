import subprocess
from datetime import datetime

namespaces = [""]

def get_namespace_uid(namespace):
    try:
        result = subprocess.run(['kubectl', 'get', 'namespace', namespace, '-o', 'jsonpath={.metadata.uid}'], capture_output=True, text=True, check=True)
        namespace_uid = result.stdout.strip()
        return namespace_uid
    except subprocess.CalledProcessError as e:
        print(f"Error getting UID for namespace {namespace}: {e}")
        return None

def create_secret(namespace, secret_name, secret_data):
    try:
        subprocess.run(['kubectl', 'create', 'secret', 'generic', secret_name, '--from-literal=token='+secret_data, '--namespace='+namespace], check=True)
        print(f"Secret {secret_name} created in namespace {namespace}")
    except subprocess.CalledProcessError as e:
        print(f"Error creating secret {secret_name} in namespace {namespace}: {e}")

def apply_commands_in_each_namespace(namespaces):
    for namespace in namespaces:
        namespace_uid = get_namespace_uid(namespace)
        if namespace_uid:
            current_date = datetime.now()
            secret_name = f"{namespace}-token-{current_date.strftime('%m-%Y')}"
            secret_data = f"some-random-token-{namespace_uid}"
            create_secret(namespace, secret_name, secret_data)

if __name__ == "__main__":
    # List of Kubernetes namespaces
    kubernetes_namespaces = ['fa', 'fl', 'lo', 'mr', 'w', 'ds']

    apply_commands_in_each_namespace(kubernetes_namespaces)

