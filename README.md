# Go-Serverless-AWSLambda

This is the basic tutorial to showcase the Simple Users CRUD operations using Go Lang + AWS Lambda functions + Dynamo DB

Below are the steps for manual deployment of the code onto the AWS Lambda functions

# Step 1: Go to Lambda Service on AWS Console 
 - Select apropriate function name.
 - Change the default execution role to Create a new role form AWS policy templete
 - Select appropriate Role name and choose Simple Microservice permissions from Dropdown and CREATE
 - Move onto Lambda details and click on Edit Runtime Settings
 - Under Handler section write name of your zip file: In my case it will be main
 - Click "UPLOAD FROM" from source code section and upload the zip file 
 
# Step 2: Go to Dynamo DB on AWS Console 
  - Create Table and provide the same name provided inside code.
  - Add partition key as primary key: My case its "email"
  - Create Table
 
# Step 3: Go to API Gateway on AWS Console 
  - Select REST API and click Build
  - Add Api Name and Create 
  - Under Actions -> Create Method -> Any
  - Select Integration type and Lambda and check Lambda Proxy Integration
  - Select the created Lamdba function
  - Save
  - Actions-> Deploy API -> Stage Name("Staging") -> Deploy
  - Copy the invoke url and test the api's
  
  # ENJOYYYYYYYYY!!!!!
  
  

 

 
