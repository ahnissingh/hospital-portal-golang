-- Drop indexes
DROP INDEX IF EXISTS idx_patients_age;
DROP INDEX IF EXISTS idx_patients_gender;
DROP INDEX IF EXISTS idx_patients_name;

-- Drop patients table
DROP TABLE IF EXISTS patients;

-- Drop users table
DROP TABLE IF EXISTS users;
