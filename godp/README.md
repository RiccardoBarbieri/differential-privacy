# Run the DP Pipeline

## Prerequisites

1. Ensure you have Go SDK installed on your system. You can download it from [golang.org](https://golang.org/dl/).
2. Make sure you have set up your Go workspace correctly.
3. Install make to use the Makefile build and run commands 
4. The current build system assumes linux OS distro

## Steps to Run the Program

1. Clone the repository:
    ```bash
    git clone https://github.com/RiccardoBarbieri/differential-privacy
    cd differential-privacy/godp
    ```
2. Build the program executable:
    ```bash
   make build
    ```
3. Run
   ```bash
   make run
   ```
   to execute the pipeline

[//]: # (### Troubleshooting)

[//]: # ()
[//]: # (If you encounter any dependency issues, try running:)

[//]: # (go mod tidy)

[//]: # ()
[//]: # ()
[//]: # (Make sure all required data files specified through parameters exist.)

[//]: # ()
[//]: # ()
[//]: # (If you face permission issues, ensure you have the necessary rights to read/write in the project directory and files.)
