-- Seed data for development/testing
-- Run this after init.sql

-- Sample events for the test user
DO $$
DECLARE
    test_user_id UUID := '00000000-0000-0000-0000-000000000001';
    base_time TIMESTAMPTZ := NOW() - INTERVAL '7 days';
BEGIN
    -- Insert sample sleep data (last 7 days)
    FOR i IN 0..6 LOOP
        INSERT INTO events (time, user_id, event_type, source, data, confidence)
        VALUES (
            base_time + (i || ' days')::INTERVAL,
            test_user_id,
            'garmin_sleep',
            'garmin',
            jsonb_build_object(
                'duration_minutes', 420 + (random() * 60)::int,
                'deep_sleep_minutes', 80 + (random() * 30)::int,
                'light_sleep_minutes', 240 + (random() * 40)::int,
                'rem_sleep_minutes', 90 + (random() * 20)::int,
                'awake_minutes', (random() * 20)::int,
                'sleep_score', 70 + (random() * 25)::int,
                'hrv_avg', 55 + (random() * 20)
            ),
            0.95
        );
    END LOOP;

    -- Insert sample activity data
    FOR i IN 0..6 LOOP
        INSERT INTO events (time, user_id, event_type, source, data, confidence)
        VALUES (
            base_time + (i || ' days')::INTERVAL + INTERVAL '10 hours',
            test_user_id,
            'garmin_activity',
            'garmin',
            jsonb_build_object(
                'activity_type', (ARRAY['strength_training', 'running', 'cycling', 'walking'])[1 + floor(random() * 4)::int],
                'duration_minutes', 30 + (random() * 60)::int,
                'calories', 200 + (random() * 300)::int,
                'avg_hr', 120 + (random() * 40)::int,
                'max_hr', 160 + (random() * 30)::int
            ),
            0.95
        );
    END LOOP;

    -- Insert sample subjective feelings (morning)
    FOR i IN 0..6 LOOP
        INSERT INTO events (time, user_id, event_type, source, data)
        VALUES (
            base_time + (i || ' days')::INTERVAL + INTERVAL '8 hours',
            test_user_id,
            'subjective_feeling',
            'manual',
            jsonb_build_object(
                'energy', 5 + (random() * 5)::int,
                'mood', 6 + (random() * 4)::int,
                'focus', 5 + (random() * 5)::int,
                'physical', 6 + (random() * 4)::int,
                'notes', 'Morning check-in'
            )
        );
    END LOOP;

    -- Insert sample meals
    FOR i IN 0..6 LOOP
        -- Breakfast
        INSERT INTO events (time, user_id, event_type, source, data, confidence)
        VALUES (
            base_time + (i || ' days')::INTERVAL + INTERVAL '8 hours 30 minutes',
            test_user_id,
            'meal',
            'llm',
            jsonb_build_object(
                'meal_type', 'breakfast',
                'photo_url', 's3://health-photos/test-breakfast.jpg',
                'macros', jsonb_build_object(
                    'calories', 400 + (random() * 200)::int,
                    'protein_g', 20 + (random() * 15),
                    'carbs_g', 50 + (random() * 30),
                    'fat_g', 15 + (random() * 10),
                    'fiber_g', 5 + (random() * 5)
                ),
                'manually_verified', false
            ),
            0.75
        );

        -- Lunch
        INSERT INTO events (time, user_id, event_type, source, data, confidence)
        VALUES (
            base_time + (i || ' days')::INTERVAL + INTERVAL '13 hours',
            test_user_id,
            'meal',
            'llm',
            jsonb_build_object(
                'meal_type', 'lunch',
                'photo_url', 's3://health-photos/test-lunch.jpg',
                'macros', jsonb_build_object(
                    'calories', 600 + (random() * 250)::int,
                    'protein_g', 35 + (random() * 20),
                    'carbs_g', 60 + (random() * 30),
                    'fat_g', 20 + (random() * 15),
                    'fiber_g', 8 + (random() * 7)
                ),
                'manually_verified', false
            ),
            0.78
        );

        -- Dinner
        INSERT INTO events (time, user_id, event_type, source, data, confidence)
        VALUES (
            base_time + (i || ' days')::INTERVAL + INTERVAL '19 hours',
            test_user_id,
            'meal',
            'llm',
            jsonb_build_object(
                'meal_type', 'dinner',
                'photo_url', 's3://health-photos/test-dinner.jpg',
                'macros', jsonb_build_object(
                    'calories', 650 + (random() * 300)::int,
                    'protein_g', 40 + (random() * 25),
                    'carbs_g', 70 + (random() * 35),
                    'fat_g', 25 + (random() * 15),
                    'fiber_g', 10 + (random() * 8)
                ),
                'manually_verified', false
            ),
            0.72
        );
    END LOOP;

    -- Insert sample supplements
    FOR i IN 0..6 LOOP
        -- Morning supplements
        INSERT INTO events (time, user_id, event_type, source, data)
        VALUES (
            base_time + (i || ' days')::INTERVAL + INTERVAL '8 hours',
            test_user_id,
            'supplement',
            'manual',
            jsonb_build_object(
                'name', 'Creatine Monohydrate',
                'dosage', '5g',
                'taken', (random() > 0.1),
                'scheduled_time', '08:00',
                'actual_time', (base_time + (i || ' days')::INTERVAL + INTERVAL '8 hours')::text
            )
        );

        -- Post-workout protein
        IF i % 2 = 0 THEN
            INSERT INTO events (time, user_id, event_type, source, data)
            VALUES (
                base_time + (i || ' days')::INTERVAL + INTERVAL '11 hours',
                test_user_id,
                'supplement',
                'manual',
                jsonb_build_object(
                    'name', 'Whey Protein',
                    'dosage', '25g',
                    'taken', true,
                    'scheduled_time', '11:00',
                    'actual_time', (base_time + (i || ' days')::INTERVAL + INTERVAL '11 hours')::text
                )
            );
        END IF;
    END LOOP;

    RAISE NOTICE 'Seed data inserted successfully for test user';
    RAISE NOTICE '- 7 days of sleep data';
    RAISE NOTICE '- 7 days of activity data';
    RAISE NOTICE '- 7 days of subjective feelings';
    RAISE NOTICE '- 21 meals (breakfast, lunch, dinner)';
    RAISE NOTICE '- Supplement logs';
END $$;

-- Sample experiment
INSERT INTO experiments (
    id,
    user_id,
    name,
    hypothesis,
    status,
    intervention,
    duration_days,
    created_at
)
VALUES (
    gen_random_uuid(),
    '00000000-0000-0000-0000-000000000001',
    'Creatine Effect on Recovery',
    '5g daily creatine monohydrate improves HRV recovery time and reduces muscle soreness after workouts',
    'proposed',
    '{"supplement": "creatine_monohydrate", "dosage": "5g", "timing": "morning"}'::jsonb,
    28,
    NOW()
);

-- Success
SELECT 'Seed data loaded!' as status;
