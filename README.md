dotf-go - Dotfiles handler written in go
-------------------------------------
Primary intended function is to keep local dotfiles and remote dotfiles up-to-date
across two or more systems. Long-term goal is to also be able to maintain
installed applications.

Dependencies
------------
go, gotk, gcc, libappindicator-gtk3

Functional requirements
-----------------------
- [ ] 1. Move a file or directory into a dotfiles dir and replace the file with a symbolic link pointing to this location
- [x] 2. When new files are pushed to the remote from sys A, they should immediately be downloaded down into the dotfiles dir of sys B
- [x] 3. When a file is added locally to the dotfiles dir, it should be uploaded as soon as possible to the remote
- [ ] 4. It should be possible to configure the settings of the application
- [x] 5. If the added files cannot be uploaded when they are added, an attempt to upload them should be made every X time.
- [x] 6. The application must then run in the background
- [x] 7. A systray icon should be visible to give status info and signal that the service is running.
