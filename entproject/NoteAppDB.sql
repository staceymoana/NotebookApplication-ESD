CREATE TABLE "User" (
	UserID INT PRIMARY KEY,
	GivenName VARCHAR(30),
	FamilyName VARCHAR(30),
	Password VARCHAR(30)
);

CREATE TABLE Note (
	NoteID INT PRIMARY KEY,
	UserID INT,
	Title VARCHAR(30),
	Contents VARCHAR(300),
	DateCreated DATE,
	DateUpdated DATE,
	FOREIGN KEY (UserID) REFERENCES "User"(UserID)
);

CREATE TABLE NoteAccess (
	NoteAccessID INT PRIMARY KEY,
	NoteID INT,
	UserID INT,
	Read BOOL,
	Write BOOL,
	FOREIGN KEY (NoteID) REFERENCES Note(NoteID),
	FOREIGN KEY (UserID) REFERENCES "User"(UserID)
);

INSERT INTO "User" VALUES 
    (1,'Ezra','Adkins','erat'),
    (2,'Kasper','Richard','in'),
    (3,'Mason','Bush','aliquet'),
    (4,'Jerry','Martinez','quisque'),
    (5,'Carla','Petersen','ultrices'),
    (6,'Deacon','Rios','sed'),
    (7,'Odysseus','Pickett','nunc'),
    (8,'Quintessa','Lee','orci'),
    (9,'Sophia','Thornton','risus'),
    (10,'Chiquita','Bass','nulla');

INSERT INTO Note VALUES 
    (1, 1, 'This is a title.', 'Contents of the first note', DATE('now'), DATE('now')),
    (2, 1, 'Second title.', 'Contents of the second note', DATE('now'), DATE('now')),
    (3, 3, 'Third title.', 'Some contents of a note blah blah', DATE('now'), DATE('now')),
    (4, 5, 'Fourth title.', 'This is some contents', DATE('now'), DATE('now')),
    (5, 10, 'Fifth title.', 'BBQ shapes', DATE('now'), DATE('now'));