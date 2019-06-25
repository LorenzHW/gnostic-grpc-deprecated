
3.0.0�
Test API for GSoC project�This is a OpenAPI description for testing my GSoC project. The name of the path defines what
will be tested and the operation object will be set accordingly.
Structure of tests:
/testParameter*   --> To test everything related to path/query parameteres
/testResponse*    --> To test everything related to respones
/testRequestBody* --> To test everything related to request bodies
others            --> Other stuff
21.0.0"�
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
�
/testResponse400StatusCodei"g*testResponse400StatusCodeBJH
400A
?
error6
4
application/json 

#/components/schemas/Error
�
!/testResponseAdditionalProperties�"�* testResponseAdditionalPropertiesB_]
200V
T
successful operation<
:
application/json&
$
"�object�

�integer�int32
s
/testResponseComponentReferenceP"N*testResponseComponentReferenceB,*
200#!
#/components/responses/Response*�
�
^
ErrorU
S�code�message�object�6

code
�integer�int32

message
	�string
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