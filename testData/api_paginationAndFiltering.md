<apidoc group="General" title="Pagination and Filtering" priority="140" />

Pagination and Filtering
-------

<p>Responses are paginated with a default page size of 20. If more results are available, the response will include "hasMore": true.</p>
<p>To retrieve the next page, use the "lastId" parameter:</p>
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" "https://[host]/api/model/customer?lastId=23"</code>  

<p>To adjust the number of results per page, use the "limit" parameter (allowed range: 20 - 100).</p>

<p>By default the responses don't include null fields. If you want to include them add includeNulls=true</p>
<code>https://[host]/api/model/customer?includeNulls=true</code>      

<h3>Field Selection</h3>
<p>You can specify which fields to include in the response. Some fields, like "id", are always included.</p>
<code>https://[host]/api/model/customer?fields=["name", "avatar"]</code>

<p>URL parameters are shown unencoded for clarity but must be encoded:</p>
<code>https://[host]/api/model/customer?fields=%5B%22name%22%2C%20%22avatar%22%5D</code>

<h3>Filtering</h3>
<p>To filter results, use the format ["field", "operator", "value"]. Examples:</p>

<h4>Filter by date</h4>
<code>https://[host]/api/model/customer?search=["created",">","05-07-2023"]</code>

<h4>Filter by null values</h4>
<code>["end","=",null]</code>

<h4>Include deleted records</h4>
<code>["deleted","in",[0,1]]</code>  

<h4>Nested filters</h4>
<code>["OR", ["status","in",[3,4]], ["price","=",0]]</code>

<h4>Filter records created or updated since a date</h4>
<code>["OR", ["created",">","2023-07-01"], ["updated",">","2023-07-01"]]</code>