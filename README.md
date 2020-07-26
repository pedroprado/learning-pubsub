**Proof of Concept: Google Cloud Platfomr Pubsub**

This project was conceived to test the utilization of GCP Pubsub
It contains three services: Pubsub Emulator, Publisher and Listener
The Publiser-api can be used to publish a message (mocked) to Pubsub:
     --send a GET request to *localhost:5000/message/publish
The Listener-api keeps polling messages from the Pubsub topic, each 3 seconds and prints it out to the console

---

## Running locally

sudo docker-compose up --build