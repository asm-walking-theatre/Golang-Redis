Golang-Redis project Citymall
Submitted by: Akash Mudhol
Methodology:
GEOADD key{maps} Longitude Latitude Member{Pincode}
Pincode Versioning:
The above is the Redis function we are using to add a particular pincode to our database. In this function the member acts as a unique “Id” like Primary Key in SQL databases. So, when a point: Point 1 is initially added for say ”431513” it gets added but when say Point 2: which is outside 1 km range of Point 1 but still has pincode 431513 we were unable to add. So, I created a Pincode versioning function that maps the previous versions and add the proper version for our requested pincode to cache. So, we can add that co-ordinate as well and every pincode entry in our cache is unique. This versioning is internal implementation of the code and won’t be visible to Client, only pincode is generated.

Edge Cases and proposed solutions:
1. The points at the boundary of 2 different pincodes.
Sol. When we receive the output via GEORADIUS command and if we get the output as 2 different pincodes then there are two solutions to it: 
a.	If cost saving is aim, we can output the pincode with highest count
b.	Calling API whenever we have 2 different pincodes as output. I personally prefer method 2 as it will eventually increase the accuracy of our Redis cache.
2.Chain effect can happen if requests of consecutive pincodes are also consecutively.
Sol. I have added the Request_Count variable in our server so we need to decide threshold say 20, i.e., after every 20 requests once we will call the API
3. We need to add latitude and longitude bounds of India so as the output pincode is of India only.
4. Inability to add same pincode is the flaw of Redis, but it is an implementation issue and I have addressed that in my code. 

File details
Link to GitHub repo: https://github.com/asm-walking-theatre/Golang-Redis
I have implemented the proposed solution and tested it for a variety of range of latitude and longitudes in India. This is curated as of now for Citymall assuming all incoming client requests are from India and we need a 6-digit pincode.
Files:
1.	main.go : This is the server file which needs to execute first and the server is started
2.	client.go : This is the client file to send requests to server and it returns pincode
Note: A sample output response of both main.go and client.go are attached in the repo
Note: Redis server needs to be running for the code to get executed.
