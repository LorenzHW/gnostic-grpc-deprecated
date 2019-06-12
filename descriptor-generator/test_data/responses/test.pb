
3.0.0Ä
Test API for GSoC projectŸThis is a OpenAPI description for testing my GSoC project. The name of the path defines what
will be tested and the operation object will be set accordingly.
Structure of tests:
/testParameter*   --> To test everything related to path/query parameteres
/testResponse*    --> To test everything related to respones
/testRequestBody* --> To test everything related to request bodies
others            --> Other stuff
21.0.0"é
g
/testResponseNativeP"N*testResponseNativeB86
200/
-
succes#
!
application/json

	Êstring
˜
/testResponseReference~"|*testResponseReferenceBca
200Z
X
successM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson
‹
/testResponseMultipleContentê"ç*testResponseMultipleContentBÇ°
200¨
¥
success™
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
ƒ
/testResponseComponentReference`"^*testResponseComponentReferenceB<:
20031
/#/components/responses/ComponentExampleResponse*Ó
Ö
Ó
ComponentExampleObjectPerson²
¯ºnameº	photoUrlsÊobjectú

id
Êintegeršint64

age
Êintegeršint64

name
:Peter
Êstring
5
	photoUrls(
&*
photoUrl(Êarrayò

	Êstringx
v
ComponentExampleResponseZ
X
successM
K
application/json7
53
1#/components/schemas/ComponentExampleObjectPerson