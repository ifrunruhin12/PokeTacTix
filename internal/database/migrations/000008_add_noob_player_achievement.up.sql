-- Add "Noob Player" achievement for losing 10 consecutive battles
INSERT INTO achievements (name, description, icon, requirement_type, requirement_value) VALUES
    ('Noob Player', 'Lose 10 battles in a row (any mode)', '😅', 'consecutive_losses', 10);
