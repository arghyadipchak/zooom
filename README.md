<h1 align="center">Zooom</h3>

A simple script to join your regular Zoom meetings faster, specially useful for students who have a lot of links to manage for their daily Zoom classes

## Getting Started

1. [Download the latest release](https://github.com/arghyadipchak/zooom/releases/latest) and untar/unzip it
2. Create a `config.json` file, an example file (`config.example.json`) is already provided
3. Now to join Zoom meetings you just need to run the `zooom` executable

## Configuration File

The `config.json` file contains the configurations for Zooom. Following keys are required:

- `sources`: Array of sources. Each source is either the location of file containing an array of Meeting objects or a http(s) endpoint that return an array of Meeting objects as the body
- `buffer`:
  - `start`: Duration before the starting time of a meeting you can join it
  - `end`: Duration after the ending time of a meeting you can join it

Example:
```json
{
  "sources": [
    "meetings.json"
  ],
  "buffer": {
    "start": "00:10",
    "end"  : "00:00"
  }
}
```
*Time & Duration are in HH:MM format*

Configuration file path can be set using `ZOOOM_CONFIG` environment variable. *(Requires >= v1.0.1)*

Example: `ZOOOM_CONFIG="~/.config/zooom.json"`

## Meeting Object

JSON object with the following keys:

- `name`: Meeting Name
- `days`: Array of the week-days (just the first 3 letters) the meeting is held
- `start`: Starting Time of the Meeting
- `end`: Ending Time of the Meeting
- `metno`: Meeting Number
- `paswd`: Meeting Password

Example:
```json
[
  {
    "name" : "Meeting1",
    "days" : ["Mon", "Tue", "Wed", "Sun"],
    "start": "09:00",
    "end"  : "12:00",
    "metno": "0000000000",
    "paswd": "999999"
  },
  {
    "name" : "Meeting2",
    "days" : ["Thu", "Fri", "Sat"],
    "start": "11:30",
    "end"  : "15:15",
    "metno": "0000000001",
    "paswd": "999998"
  }
]
```
