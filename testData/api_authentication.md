<apidoc group="General" title="Authentication" priority="110" />

     
Authentication
-----------

<p>Include an HTTP header named "key" with your API Key to authenticate requests.</p>
<p>You also need to specify a "tenant" header.</p>

<p>Example:</p>
<code>curl -H "key: [apiKey]" -H "tenant: [tenant]" "https://[host]/api/model/customer"</code>
