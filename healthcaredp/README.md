# Run the DP Pipeline

## Prerequisites

1. Ensure you have Go SDK installed on your system. You can download it from [golang.org](https://golang.org/dl/).
2. Make sure you have set up your Go workspace correctly.

## Steps to Run the Program

1. Clone the repository:
    ```bash
    git clone https://github.com/RiccardoBarbieri/differential-privacy
    cd differential-privacy/healthcaredp
    ```
2. Install dependencies:
    ```bash
    go mod download
    ```
3. Navigate to the main directory:
    ```bash
    cd main
    ```
4. Build the program executable:
    ```bash 
    go build -o build/healthcaredp .
    ```
5. Run
   ```bash
   ./build/healthcaredp -h
   ```
   to show documentation for the parameters.

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
