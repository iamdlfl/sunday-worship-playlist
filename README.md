# Sunday Worship Playlist
The purpose of this script is to update a Spotify playlist with the songs that will be sung during the upcoming Sunday service for Creekside Community Church. 
This is done with the following steps:
1. Get songs for the upcoming Sunday using Planning Center Online API.
2. Get the "Upcoming Sunday" playlist from Spotify.
3. Replace the songs on this playlist with the songs retrieved from Planning Center.
4. Create a new Spotify playlist for the Sunday's date with these same songs, to act as a historical record.
5. Send an email with success or failure.

## Configuration
This project uses an INI settings file with sections for the mailer, Spotify and PCO. The relevant auth information is contained in this file. 

## TODO
- Update code to perform steps 2 & 3. Current version just does 1, 4 and 5, but having a single "Upcoming Sunday" playlist means the link won't have to be updated in the weekly email. 
