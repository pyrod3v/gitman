import os
import shutil
import platform
import sys
import subprocess

if __name__ == "__main__":
    system = platform.system()
    if system in ("Linux", "Darwin"):
        dest_directory = "/usr/local/bin"
        path = "gitman"
    elif system == "Windows":
        dest_directory = os.path.join(os.environ.get("ProgramFiles", "C:\\Program Files"), "bin")
        path = "gitman.exe"
    else:
        print(f"Unsupported operating system: {system}")
        sys.exit(1)

    if len(sys.argv) > 1:
        path = sys.argv[1]

    if not os.path.exists(path):
        try:
            subprocess.run(["go", "build", "-o", path], check=True)
            print(f"Successfully built binary: {path}")
        except subprocess.CalledProcessError as e:
            print(f"Error building binary: {e}")
            sys.exit(1)

    os.makedirs(dest_directory, exist_ok=True)
    dest = os.path.join(dest_directory, path)
        
    try:
        shutil.move(path, dest)
        print(f"Successfully moved to {dest}")
    except Exception as e:
        print(f"Error moving file: {e}")
        sys.exit(1)
