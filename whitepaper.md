[![Build Status](https://travis-ci.org/alexshemesh/claptrap.svg?branch=master)](https://travis-ci.org/alexshemesh/claptrap)
### Whitepaper
Cryptocurrency world is diverse and hard to grasp.But pays up nicely.
There are tons of sites, exchanges, sources of information that had to be monitored.
I have full time job that wont allow me to spend more than 1 hour a day to manage my crypto assets.
Telegram bot can be very useful to combine and present information in secure manner.
And writing software is a good way to learn ecosystem.

# Goals:
1. Accumulate information from multiple sources to provide customizable reports
  *  I have multiple accounts on different exchanges, its takes a lot of time to visit all of them and compile some kind of report to understand where do i stand.
2. Implement funding bot functionality to lend cryptocurrency
  *  Currently i use cryptolend.net. Its awesome no doubt but its does not look so hard to implement. So why trust 3rd party service with my money?
3. Accumulate statistical data for later analyses.
  *  There is just too much information to process. 
  *  For example arbitrage between exchanges looks tempting yet its very hard to estimate if it will succeed because of various reasons. 
  *  Statistical data on length of transactions, error rates, perhaps some estimations for rate volatility can be very usefull in implmenting full fledged trading system.


# Implementation:
1. [Telegram](https://telegram.org/) serves as a queue of requests. That queue can be easily distributed in the background.
2. [Telegram](https://telegram.org/) provides security for communication. No need to provide any inbound connection.
3. We will use [Go](https://golang.org/) as primary development language
  [Go](https://golang.org/) allows fast development and deployment of the project. No need to worry about environment or os type.
4. We will develop code in TDD manner. To make it more stable and predictable, and to achieve better architectual decisions.
5. All sensitive information will be stored in Hashicorp vault server. 
 * Its easy to manage
 * Very fast
 * Very stable
 * Scalable
 * Has flexible user management system.
6. We will use [Travis](https://travis-ci.org/) as CI service. 
7. I didnt decided about DB server yet. 
  *  I usualy use [MySQL](https://www.mysql.com/) and [MongoDB](https://www.mongodb.com/) in most of my projects. 
  *  Another candidate can be [nuodb](https://www.nuodb.com/). Looks scalable and promising.
  *  [Vault](https://www.vaultproject.io/) has option to encrypt data so only authenticated user can decrypt it.  I want to use it to encrypt individual records to provide more security.
8. I will use [ELK](https://www.elastic.co/products) stack to accumulate statistics about system functionality
9. Service will be deployed in AWS at first.
  * [Terraform](https://www.terraform.io/) + [Ansible](https://www.ansible.com/) provisioning
  * Heavy usage of dockers
  * [DataDog](https://www.datadoghq.com/) for monitoring

# Command line interface
https://github.com/alexshemesh/claptrap/wiki/CLI-docs


