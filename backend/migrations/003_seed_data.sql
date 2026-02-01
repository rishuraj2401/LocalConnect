-- Seed Data for LocalConnect
-- This file creates demo users, profiles, reviews, and other data

-- Insert demo clients
INSERT INTO users (name, email, phone, password_hash, role) VALUES
('Sarah Johnson', 'sarah@example.com', '+1-555-0101', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'client'),
('Michael Chen', 'michael@example.com', '+1-555-0102', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'client'),
('Emma Davis', 'emma@example.com', '+1-555-0103', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'client'),
('James Wilson', 'james@example.com', '+1-555-0104', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'client')
ON CONFLICT (email) DO NOTHING;

-- Insert demo workers
INSERT INTO users (name, email, phone, password_hash, role) VALUES
('Robert Martinez', 'robert@example.com', '+1-555-0201', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Jennifer Lee', 'jennifer@example.com', '+1-555-0202', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('David Anderson', 'david@example.com', '+1-555-0203', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Maria Garcia', 'maria@example.com', '+1-555-0204', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Thomas Brown', 'thomas@example.com', '+1-555-0205', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Lisa Taylor', 'lisa@example.com', '+1-555-0206', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Kevin White', 'kevin@example.com', '+1-555-0207', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Anna Rodriguez', 'anna@example.com', '+1-555-0208', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Christopher Harris', 'chris@example.com', '+1-555-0209', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Patricia Clark', 'patricia@example.com', '+1-555-0210', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Daniel Lewis', 'daniel@example.com', '+1-555-0211', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker'),
('Jessica Martinez', 'jessica@example.com', '+1-555-0212', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'worker')
ON CONFLICT (email) DO NOTHING;

-- Create worker profiles with realistic data
INSERT INTO worker_profiles (user_id, category_id, location, rate, experience_years, bio) 
SELECT 
    u.id,
    c.id,
    location,
    rate,
    experience,
    bio
FROM (
    SELECT 'robert@example.com' as email, 'carpenter' as category, 'New York, NY' as location, 150 as rate, 8 as experience, 'Experienced carpenter specializing in custom furniture, cabinetry, and home renovations. Licensed and insured with a focus on quality craftsmanship and attention to detail.' as bio
    UNION ALL
    SELECT 'jennifer@example.com', 'painter', 'Los Angeles, CA', 120, 6, 'Professional painter with expertise in interior and exterior painting. Skilled in color consultation, wallpaper removal, and various painting techniques including faux finishes and murals.' as bio
    UNION ALL
    SELECT 'david@example.com', 'labour', 'Chicago, IL', 100, 5, 'Reliable general laborer available for various projects including loading/unloading, demolition, cleanup, and general construction assistance. Strong work ethic and physically fit.' as bio
    UNION ALL
    SELECT 'maria@example.com', 'cook', 'Houston, TX', 130, 10, 'Professional chef with 10 years experience in various cuisines including Mexican, Italian, and American. Available for private events, meal prep, and cooking classes. Food handler certified.' as bio
    UNION ALL
    SELECT 'thomas@example.com', 'carpenter', 'Phoenix, AZ', 140, 7, 'Skilled carpenter specializing in deck building, framing, and finish carpentry. Expert in reading blueprints and working with various wood types. Always on time and professional.' as bio
    UNION ALL
    SELECT 'lisa@example.com', 'home tution', 'Philadelphia, PA', 80, 12, 'Certified teacher with Masters in Education. Offering tutoring services for K-12 students in Math, Science, and English. Patient, encouraging teaching style with proven results.' as bio
    UNION ALL
    SELECT 'kevin@example.com', 'teacher', 'San Antonio, TX', 90, 9, 'Experienced educator specializing in high school mathematics and physics. Available for one-on-one tutoring and small group instruction. SAT/ACT prep specialist.' as bio
    UNION ALL
    SELECT 'anna@example.com', 'househelp', 'San Diego, CA', 110, 4, 'Trustworthy and detail-oriented housekeeper with excellent references. Services include deep cleaning, organizing, laundry, and light cooking. Background checked and insured.' as bio
    UNION ALL
    SELECT 'chris@example.com', 'painter', 'Dallas, TX', 125, 11, 'Master painter with over a decade of experience. Specializing in residential and commercial projects. Expert in preparation, color matching, and achieving perfect finishes.' as bio
    UNION ALL
    SELECT 'patricia@example.com', 'cook', 'San Jose, CA', 140, 8, 'Private chef and caterer specializing in healthy meal preparation and dietary restrictions. Experienced with vegan, keto, and gluten-free cooking. ServSafe certified.' as bio
    UNION ALL
    SELECT 'daniel@example.com', 'labour', 'Austin, TX', 95, 3, 'Hardworking and dependable laborer seeking opportunities in construction, moving, and landscaping. Quick learner with strong physical stamina and positive attitude.' as bio
    UNION ALL
    SELECT 'jessica@example.com', 'househelp', 'Jacksonville, FL', 105, 7, 'Professional house cleaner with 7 years experience. Thorough, reliable, and respectful of your home. Eco-friendly cleaning products available. Flexible scheduling.' as bio
) AS demo
JOIN users u ON u.email = demo.email
JOIN categories c ON c.name = demo.category
ON CONFLICT (user_id) DO NOTHING;

-- Add reviews for profiles
DO $$
DECLARE
    profile_record RECORD;
    client_id UUID;
    review_count INT;
BEGIN
    -- Get a client user ID
    SELECT id INTO client_id FROM users WHERE role = 'client' LIMIT 1;
    
    -- Add reviews for each profile
    FOR profile_record IN (SELECT id FROM worker_profiles) LOOP
        -- Add 2-4 random reviews per profile
        review_count := 2 + floor(random() * 3)::int;
        
        FOR i IN 1..review_count LOOP
            INSERT INTO reviews (profile_id, user_id, rating, comment)
            VALUES (
                profile_record.id,
                client_id,
                4 + floor(random() * 2)::int, -- Random rating between 4-5
                CASE floor(random() * 5)::int
                    WHEN 0 THEN 'Excellent work! Very professional and completed the job on time. Would highly recommend!'
                    WHEN 1 THEN 'Great experience! Quality work and fair pricing. Will definitely hire again.'
                    WHEN 2 THEN 'Very satisfied with the service. Punctual, skilled, and pleasant to work with.'
                    WHEN 3 THEN 'Outstanding results! Exceeded my expectations in every way. Five stars!'
                    ELSE 'Reliable and skilled professional. Great attention to detail and customer service.'
                END
            );
        END LOOP;
    END LOOP;
END $$;

-- Update profile stats based on reviews
UPDATE worker_profiles wp
SET 
    review_count = (SELECT COUNT(*) FROM reviews WHERE profile_id = wp.id),
    average_rating = (SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE profile_id = wp.id)
WHERE EXISTS (SELECT 1 FROM reviews WHERE profile_id = wp.id);

-- Add some upvotes
DO $$
DECLARE
    profile_record RECORD;
    client_record RECORD;
    upvote_count INT;
BEGIN
    FOR profile_record IN (SELECT id FROM worker_profiles) LOOP
        upvote_count := 3 + floor(random() * 8)::int; -- 3-10 upvotes
        
        FOR i IN 1..upvote_count LOOP
            -- Get a random client
            SELECT id INTO client_record FROM users WHERE role = 'client' ORDER BY random() LIMIT 1;
            
            INSERT INTO upvotes (profile_id, user_id)
            VALUES (profile_record.id, client_record.id)
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- Update upvote counts
UPDATE worker_profiles wp
SET upvote_count = (
    SELECT COUNT(*) FROM upvotes WHERE profile_id = wp.id
);

-- Add some contact requests
DO $$
DECLARE
    profile_record RECORD;
    client_record RECORD;
BEGIN
    FOR profile_record IN (SELECT id, user_id FROM worker_profiles ORDER BY random() LIMIT 6) LOOP
        -- Get a random client
        SELECT id INTO client_record FROM users WHERE role = 'client' ORDER BY random() LIMIT 1;
        
        INSERT INTO contact_requests (profile_id, user_id, message, phone_shared)
        VALUES (
            profile_record.id,
            client_record.id,
            CASE floor(random() * 4)::int
                WHEN 0 THEN 'Hi! I need help with a project next week. Are you available?'
                WHEN 1 THEN 'I would like to get a quote for my upcoming project. Please contact me.'
                WHEN 2 THEN 'Your profile looks great! I have a job that might interest you. Can we discuss?'
                ELSE 'I am looking for your services. Please let me know your availability.'
            END,
            true
        );
    END LOOP;
END $$;

COMMIT;
