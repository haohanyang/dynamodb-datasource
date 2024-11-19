import urllib.request
import json
import tempfile
import zipfile
import os
import glob
import subprocess
import shutil
from sys import platform

dist_dir = "dynamodb-datasource"
if os.path.isdir(dist_dir):
    shutil.rmtree(dist_dir)

# Download the latest release build
with urllib.request.urlopen(
    "https://api.github.com/repos/haohanyang/dynamodb-datasource/releases/latest"
) as response:
    data = json.loads(response.read())

    for asset in data["assets"]:
        if asset["content_type"] == "application/zip":
            with tempfile.NamedTemporaryFile() as temp_file:
                print("Downloading " + asset["name"] + "...")
                urllib.request.urlretrieve(
                    asset["browser_download_url"], temp_file.name
                )
                print("Extracting files to " + dist_dir)
                with zipfile.ZipFile(temp_file.name, "r") as zip_ref:
                    zip_ref.extractall(os.getcwd())
    os.rename("haohanyang-dynamodb-datasource", "dynamodb-datasource")

# Grant execute permission go binaries
for bin in glob.glob("dynamodb-datasource/gpx_dynamodb_datasource_*"):
    os.chmod(bin, os.stat(bin).st_mode | 755)

if platform == "linux" or platform == "linux2":
    # "sudo" only for linux
    subprocess.call(
        ["sudo", "docker", "compose", "-f", "docker-compose.prod.yaml", "up", "-d"]
    )
else:
    subprocess.call(["docker", "compose", "-f", "docker-compose.prod.yaml", "up", "-d"])
