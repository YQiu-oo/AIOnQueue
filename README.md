# AIOnQueue


Introduction
------------

AIONQueue is a simplified, AI-integrated message queue system inspired by Kafka. It focuses on providing a lightweight, efficient, and intelligent message handling mechanism, suitable for handling high-throughput, real-time data with smart routing and anomaly detection capabilities.

Features
--------

*   **Simple Message Queuing**: Basic produce-consume functionality.
*   **AI-Enhanced Routing**: Intelligent message distribution based on consumer load and processing capabilities.
*   **Anomaly Detection**: Real-time system monitoring and anomaly detection using AI.
*   **Scalability**: Designed to be scalable and efficient.
*   **Easy Integration**: Simplified setup and integration process.

Getting Started
---------------

### Prerequisites

*   Go version 1.15 or higher
*   \[Optional\] TensorFlow/PyTorch for AI components


### Test
Clone the repository to your local machine:

```
git clone [repository URL]
```
Then, navigate to the project directory:

```
cd AIONQueue
```
Running the RPC Server
Before running the tests, you need to start the RPC server. This can be done by executing:
```
go run .
```
Alternatively, you can build the project first and then run the built executable:
```
go build
./AIONQueue
```
Running the Tests
AIONQueue uses Go's standard testing framework for testing. To run the tests, ensure the RPC server is up and running, then execute the following command in the root directory of the project:

```
go test
```
This will execute all the test cases in the project.
