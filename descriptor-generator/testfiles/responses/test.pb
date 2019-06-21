
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
/testResponseReferenceh"f*testResponseReferenceBMK
200D
B
success7
5
application/json!

#/components/schemas/Person
�
/testResponseMultipleContent�"�*testResponseMultipleContentB��
200z
x
successm
5
application/json!

#/components/schemas/Person
4
application/xml!

#/components/schemas/Person
400
	
failure
O
/testResponse400StatusCode1"/*testResponse400StatusCodeB
400	

error
s
/testResponseComponentReferenceP"N*testResponseComponentReferenceB,*
200#!
#/components/responses/Response*�
�
�
Person�
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

	�stringR
P
ResponseD
B
success7
5
application/json!

#/components/schemas/Person