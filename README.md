**Proof of Concept: Google Cloud Platfomr Pubsub**

This project was conceived for testing the use of GCP Pubsub Docker Image
It contains 3 services: Pubsub, Publisher and Listener
The Publiser-api can be used for publishing a message (mocked) to the Pubsub:
     send a GET request to *localhost:5000/message/publish
The Listener-api keeps pulling messages from the Pubsub topic, each 3 seconds and prints it to the console

---

## Running locally

sudo docker-compose up --build