INSERT INTO "User" VALUES 
    (DEFAULT,'Ezra','Adkins','password'),
    (DEFAULT,'Kasper','Richard','password'),
    (DEFAULT,'Mason','Bush','password'),
    (DEFAULT,'Jerry','Martinez','password'),
    (DEFAULT,'Carla','Petersen','password'),
    (DEFAULT,'Deacon','Rios','password'),
    (DEFAULT,'Odysseus','Pickett','password'),
    (DEFAULT,'Quintessa','Lee','password'),
    (DEFAULT,'Sophia','Thornton','password'),
    (DEFAULT,'Chiquita','Bass','password');

INSERT INTO Note VALUES 
    (DEFAULT, 1, 'This is a title.', 'Contents of the first note', DATE('now'), DATE('now')),
    (DEFAULT, 1, 'Second title.', 'Contents of the second note', DATE('now'), DATE('now')),
    (DEFAULT, 3, 'Third title.', 'Some contents of a note blah blah', DATE('now'), DATE('now')),
    (DEFAULT, 5, 'Fourth title.', 'This is some contents', DATE('now'), DATE('now')),
    (DEFAULT, 10, 'Fifth title.', 'BBQ shapes', DATE('now'), DATE('now'));

INSERT INTO NoteAccess VALUES
	(DEFAULT, 1, 2, true, true),
	(DEFAULT, 1, 3, true, true);

INSERT INTO SharedSettings VALUES
	(DEFAULT, 1, 2, true, true, 'SharedSettings Test'),
	(DEFAULT, 1, 3, true, true, 'SharedSettings Test');