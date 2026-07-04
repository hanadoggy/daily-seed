// MongoDB seed script — runs on first docker-compose up.
// Populates sample tasks and habits for development convenience.

db = db.getSiblingDB('daily_seed');

// Only seed if the tasks collection is empty.
if (db.tasks.countDocuments() === 0) {
  print('Seeding tasks...');
  db.tasks.insertMany([
    {
      _id: 'task_001',
      section: 'japanese',
      title: 'Memorize Kanji',
      type: 'quantitative',
      metrics: { dailyTarget: 10, totalTarget: 500 },
      conditions: { weather: 'any', mode: 'any' },
      status: 'active',
    },
    {
      _id: 'task_002',
      section: 'japanese',
      title: 'Read NHK News',
      type: 'boolean',
      metrics: { dailyTarget: 1, totalTarget: 0 },
      conditions: { weather: 'any', mode: 'any' },
      status: 'active',
    },
    {
      _id: 'task_003',
      section: 'dev',
      title: 'LeetCode Problems',
      type: 'quantitative',
      metrics: { dailyTarget: 3, totalTarget: 100 },
      conditions: { weather: 'any', mode: 'any' },
      status: 'active',
    },
    {
      _id: 'task_004',
      section: 'self_dev',
      title: 'Read Book (pages)',
      type: 'quantitative',
      metrics: { dailyTarget: 20, totalTarget: 300 },
      conditions: { weather: 'any', mode: 'any' },
      status: 'active',
    },
  ]);
  print('Inserted 4 tasks.');
} else {
  print('Tasks collection already has data, skipping seed.');
}

if (db.habits.countDocuments() === 0) {
  print('Seeding habits...');
  db.habits.insertMany([
    {
      _id: 'habit_001',
      title: 'Meditate instead of using smartphone',
      category: 'mindfulness',
      status: 'active',
    },
    {
      _id: 'habit_002',
      title: 'Morning stretching routine',
      category: 'health',
      status: 'active',
    },
    {
      _id: 'habit_003',
      title: 'Write gratitude note',
      category: 'mindfulness',
      status: 'active',
    },
  ]);
  print('Inserted 3 habits.');
} else {
  print('Habits collection already has data, skipping seed.');
}

print('Seed complete.');
