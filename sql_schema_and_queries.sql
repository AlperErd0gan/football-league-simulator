-- team_models table
CREATE TABLE team_models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    strength INTEGER,
    points INTEGER,
    goals_scored INTEGER,
    goals_against INTEGER,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- match_models table
CREATE TABLE match_models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week INTEGER,
    home TEXT,
    away TEXT,
    home_goals INTEGER,
    away_goals INTEGER,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);



-- Get all teams
SELECT * FROM team_models;

-- Get current league 
SELECT name, points, goals_scored - goals_against AS gd
FROM team_models
ORDER BY points DESC, gd DESC;

-- Get all match results
SELECT * FROM match_models ORDER BY week ASC;

-- Get last played week
SELECT MAX(week) AS last_week FROM match_models;

-- Delete all data 
DELETE FROM team_models;
DELETE FROM match_models;
DELETE FROM sqlite_sequence WHERE name = 'team_models';
DELETE FROM sqlite_sequence WHERE name = 'match_models';

-- Reset auto-increment 
DELETE FROM sqlite_sequence WHERE name = 'team_models';
DELETE FROM sqlite_sequence WHERE name = 'match_models'
