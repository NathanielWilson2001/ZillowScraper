# goLangZillow

This is a fun side project I am doing to learn and better unstand Go.

The purpose is to create and send API requests to Zillow.com to gather information related to housing in a specific region. Project is still in progress, I will add more
detail here in the future. 

Requirements:

- Must have Go installed
- Must install fyne.io

Current State:
![image](https://github.com/NathanielWilson2001/ZillowScraper/assets/97745329/d02eb4bd-a45c-4494-8b24-aea067d7ca08)
- Currently the app has a basic GUI which requires the input of four GPS coordinates and a number of pages, when submitted the input is sent to the makeRequest function which makes the request and gathers the data
- Future development will be to output the results to the GUI, add data validation
