# Resource & Organization Authorization Workflow

##  Moderator Setup
When a user registers using the current workflow he has the option to create an organization,
upon creating his own organization there are 3 roles created with it initially and are linked to it: 
- User 
- Admin
- Moderator

The user who creates the organization gets assigned the moderator role by default.

The next step is to add resources to your organization. Currently, this is only textual information.

The next step is to set up permissions for each resource and what can each role do to that resource. These are the
current available permissions you can set for each role and resource:
- Read
- Write
- Delete

A moderator also has the ability to set up additional roles for his organization.

The final step in this workflow is to add users and assign them roles. 

## User Perspective 

A user with the `user` role can get the ability to get resources (currently). 

To get access to a resource you first must explicitly request access by which you gain 15 minutes access to it. 

When getting a resource there is a check to see if you have access to that resource. If you do you will proceed to it. 

If you don't an error message is thrown that tells you that you need to request access to it. 

A user also can check to which resources he currently has access to and when that access will expire. And he can also 
check to see if he has access to a certain resource.

Access to a certain resource is given if he meets the setup permissions from the moderator to that resource.