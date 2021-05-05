CREATE TABLE IF NOT EXISTS country
(
    countryId NUMERIC(4) PRIMARY KEY NOT NULL DEFAULT 0,
    countryName VARCHAR(20) NOT NULL
);
CREATE TABLE IF NOT EXISTS parcelUser
(
    userId INTEGER PRIMARY KEY,
    userEmail TEXT NOT NULL,
    userPhone VARCHAR(16) NOT NULL,
    userParcelWeight NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    countryId NUMERIC(4) NOT NULL DEFAULT 0,
    CONSTRAINT fk_user FOREIGN KEY(countryId) REFERENCES country(countryId)
);

INSERT INTO country(countryId, countryName) VALUES(0,'Unidentified');
INSERT INTO country(countryId, countryName) VALUES(1,'Cameroon');
INSERT INTO country(countryId, countryName) VALUES(2,'Ethiopia');
INSERT INTO country(countryId, countryName) VALUES(3,'Morocco');
INSERT INTO country(countryId, countryName) VALUES(4,'Mozambique');
INSERT INTO country(countryId, countryName) VALUES(5,'Uganda');

