<apidoc group="General" title="Quick start" priority="100" />


 <h2 id="intro">Introduction</h2>
     
Some examples:

#### Search available times  
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" "https://[host]/api/bookings/searchAvailability"</code>

Response:  
<code>
[
   {
      "area": 1,
      "areaName": "Golf",
      "typeName": "GF 18",
      ...
   },
   {
      "area": 1,
      "areaName": "Golf",
      "typeName": "GF 9",
      ...
   }
]
</code>


#### Get a list of customers   
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" "https://[host]/api/model/customer"</code>

Response:  
<code>
{
   "data": [
      {
         "avatar": "//upload/view/iNpLsKCheu8nw2rCErMtdVAtoY3So9bM94bordf5Bu1Z",
         "birthdate": "1980-03-10T00:00:00+01:00",
         "email": "alicia@example.com_deleted_1744044361",
         "email2": null,
         "firstName": "Alicia",
         ...
      },
      ...
   ],
   "hasMore": true,
   "lastId": 30
}
</code>

#### Retrieve a customer by id  
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" "https://[host]/api/model/customer/2"</code>

Response:  
<code>
{
    "id": 1,
    "name": "Alicia García",
    ...
}
</code>

#### Create a customer
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" -d "{"name":"Bill"}" "https://[host]/api/model/customer"</code>

Response:  
<code>
{
    "id": 2,
    "name": "John",
    ...
}
</code>

#### Update a customer  
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" -d '{"name":"Bill"}' "https://[host]/api/model/customer/2"</code>

#### Delete a customer  
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" -X DELETE "https://[host]/api/model/customer/2"</code>
