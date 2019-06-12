
3.0.0�
Test API for GSoC project�This is a OpenAPI description for testing my GSoC project. The name of the path defines what
will be tested and the operation object will be set accordingly.
Structure of tests:
/testParameter*   --> To test everything related to path/query parameteres
/testResponse*    --> To test everything related to respones
/testRequestBody* --> To test everything related to request bodies
others            --> Other stuff
21.0.0"�
g
/testResponseNativeP"N*testResponseNativeB86
200/
-
succes#
!
application/json

	�string
�
/testResponseReference~"|*testResponseReferenceBca
200Z
X
successM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson
�
/testResponseMultipleContent�"�*testResponseMultipleContentB��
200�
�
success�
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson
J
application/xml7
53
1#/components/schemas/ComponentExampleObjectPerson
400
	
failure
O
/testResponse400StatusCode1"/*testResponse400StatusCodeB
400	

error
�
/testResponseComponentReference`"^*testResponseComponentReferenceB<:
20031
/#/components/responses/ComponentExampleResponse*�
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

	�stringx
v
ComponentExampleResponseZ
X
successM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson