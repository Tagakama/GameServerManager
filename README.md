GameServerManager  is an application designed to launch game servers on local machines (with private/gray IP addresses) that are inaccessible directly from the internet. 
Description 

This utility is created for managing game servers in environments with limited external access. To function properly, GameServerManager  requires a proxy server , which is included in the repository and has been tested on dedicated Linux servers. 

In our project, this tool was used in conjunction with the Photon Fusion  framework during the development and testing of our networked game. 
Use Case 

    Local Server Without Public IP : GameServerManager runs on a local machine where game servers are hosted.
    Proxy Server with Public IP : The proxy server operates on a separate machine with a public IP, enabling external clients to connect.
    Networked Game Development : Ideal for testing network components of your game, especially when working with Photon Fusion.
     

Features 

    Private Server Management : Ability to run servers on machines without a public IP.
    Integration with Proxy Server : Facilitates access to local servers from the internet.
    Fully Commented Code : All source code is well-commented for ease of understanding.
    Ready-to-Use Proxy Server : Includes a pre-configured proxy server tested on Linux.
     

Requirements 

    Main Dependencies :
        Go version 1.18 or higher.
        Operating System: Windows/Linux.
         
    Additional :
        A proxy server with a public IP for external client access.
         
     
