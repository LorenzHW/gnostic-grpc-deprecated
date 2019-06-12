
3.0.0�
Test API for GSoC project�This is a OpenAPI description for testing my GSoC project. The name of the path defines what
will be tested and the operation object will be set accordingly.
Structure of tests:
/testParameter*   --> To test everything related to path/query parameteres
/testResponse*    --> To test everything related to respones
/testRequestBody* --> To test everything related to request bodies
others            --> Other stuff
21.0.0"�
�
/testRequestBody|"z*testRequestBody:Q
OM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPersonB
200
	
success
�
/testRequestBodyReferencen"l*testRequestBodyReference::8
6#/components/requestBodies/ComponentExampleRequestBodyB
200
	
success*�
�
�
ComponentExampleObjectPerson�
��name�	photoUrls�object��

id
�integer�int64

age
�integer�int64

name
:Peter
�string
5
	photoUrls(
&*
photoUrl(�array�

	�string*�
�
ComponentExampleRequestBodyy
w
$A JSON object containing informationM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson