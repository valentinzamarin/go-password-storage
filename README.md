# Go Password Storage

This is a small application written in Golang that I created for personal use.

Stack: 
- Go 1.23+
- Fyne (GUI)
- GORM (ORM)
- SQLite
- Argon2id (key derivation)
- NaCl/secretbox (encryption)

Architecture
- Domain-Driven Design (DDD)
- CQRS pattern
- Clean Architecture

to do: 
~~- add delete, uptade methods in domain repo ~~
~~- list update~~
- modify pass item ( in list views ) to show description ( create toggle button to show/hide )
- search field
- add db creationAt field, make a reminder about change pass 

comment: 
- cqrs is overkill, but hey, its a personal project that i constantly refactor
