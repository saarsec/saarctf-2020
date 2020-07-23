import random
import psycopg2

CONNECTION = 'host=localhost port=5432 dbname=mensaar user=mensaar password=raasnem'
conn = psycopg2.connect(CONNECTION)

c = conn.cursor()

c.execute('''CREATE TABLE "feedback" ( food INTEGER, "user" TEXT, obj TEXT, cook_token TEXT )''')

c.execute('''CREATE TABLE "user_profile" ( name TEXT, email TEXT, pwd TEXT, token TEXT, gender TEXT, ethnicity TEXT, PRIMARY KEY(email))''')

c.execute('''CREATE TABLE "food" ( id SERIAL PRIMARY KEY, name TEXT UNIQUE, ingredients TEXT, pic_src TEXT, rating INTEGER )''')

c.execute('''INSERT INTO food (name,ingredients,pic_src,rating) VALUES 
 ('Dibbelabbes','potato, beef jerky, onions, parsley, eggs, salt, pepper, maggi','pictures/food/dibbelabbes.jpg',4),
 ('Schales','potato, dried meat, leek, egg, maggi','pictures/food/schales.jpg',3),
 ('Gefillde','potato, liver sausage, butter, flour, dired meat, broth, milk, cream, maggi','pictures/food/gefillde.jpg',2),
 ('Geheirate','eggs, flour, mineral water, maggi, potato, dried meat, cream','pictures/food/geheirate.jpg',5),
 ('Hoorische','potato, breadcrumb, eggs, flour, maggi','pictures/food/hoorische.jpg',3),
 ('Grommbeerkerschdscher','potato, onion, maggi','pictures/food/krommbeerkerschdscher.jpg',4),
 ('Faasendkiechelcher','flour, butter, sugar, yeast, milk, eggs, maggi','pictures/food/faasendkiechelcher.jpg',5),
 ('Grommbeersalad','potato, vinegar, cream, maggi, dried meat, eggs','pictures/food/kartoffelsalat.jpg',4),
 ('Grommbeernsupp','potato, sour cream, sausages, butter, carrots, celery, maggi','pictures/food/kartoffelsuppe.jpg',2),
 ('Plattgeschmelzte','potato, butter, onions, maggi','pictures/food/plattgeschmeltzde.jpg',4),
 ('Bettseichersalat','hawkbit, dried meat, onion, eggs, vinegar, oil, potato, maggi','pictures/food/bettseicher.jpg',5),
 ('Schwenker','meat, curry, paprika, oil, maggi, onion','pictures/food/schwenker.jpg',5),
 ('Frikadellen','minced meat, onion, old bread, egg, paprika, garlic, maggi, butter','pictures/food/frikadelle.jpg',5),
 ('Lyoner mit Weck','meat, bun, maggi','pictures/food/lyoner.jpg',3),
 ('Merguez','lamb shoulder meat, pork, garlic, water, maggi','pictures/food/merguez.jpg',4),
 ('Bohnensupp','potato, dried meat, onions, beans, maggi, fondor','pictures/food/bohnensuppe.jpg',3),
 ('Dreggische Grommbeeren','potato, liver sausage, blood sausage, maggi','pictures/food/dreggische.jpg',3),
 ('Weck mit Eiaufstrich','bun, egg, flour, bacon, milk, chives, maggi','pictures/food/eiaufstrich.jpg',4),
 ('Kappesmängs','potato, carrot, onion, sugar, maggi','pictures/food/kappes.jpg',5),
 ('Grombeernpuffer','potato, egg, flour, applesauce, maggi','pictures/food/kartoffelpuffer.jpg',4),
 ('Kohlroulade','white cabbage, minced meat, onion, egg, parsley, mustard, maggi','pictures/food/rouladen.jpg',5),
 ('Leberknödel','ground liver, minced meat, egg, flour, breadcrumbs, marjoram, maggi, fondor','pictures/food/leberknoedel.jpg',5),
 ('Lyonerpfanne','potatoe, lyoner, onion, marjoram, maggi','pictures/food/lyonerpfanne.jpg',4),
 ('Mehlknäpp','mineral water, milk, curd, egg, flour, maggi','pictures/food/wasserspatzen.jpg',5),
 ('Querbeet','potatoe, beef, pork, oil, nutmeg, maggi','pictures/food/quer.jpg',3),
 ('Rostige Ritter','bun, milk, egg, sugar, cinnamon, custard','pictures/food/ritter.jpg',4),
 ('Sauerbraten mit Schneebällcher un Rotkraut','sauerbraten, potatoes, parsley, egg, flour, nutmeg, red cabbage, onion, apple, vinegar, juniper berries, bay leaves','pictures/food/sauerbraten.jpg',4);
''')

c.execute('''CREATE TABLE seat ( row INTEGER, tablenumber INTEGER, seatrow INTEGER, seatnumber INTEGER, reserved_by TEXT, last_reservation TIMESTAMP )''')

for a in range(1, 4):
    for b in range(1, 5):
        for i in range(1, 9):
            c.execute('''INSERT INTO seat VALUES ({:d}, {:d}, {:d}, {:d}, '')'''.format(a, b, 1, i))
            c.execute('''INSERT INTO seat VALUES ({:d}, {:d}, {:d}, {:d}, '')'''.format(a, b, 2, i))


c.execute('''CREATE TABLE "menu" ( id SERIAL PRIMARY KEY, food INTEGER, day TEXT, date TEXT )''')

# Insert first Menu
for d in ["monday", "tuesday", "wednesday", "thursday", "friday"]:
    for f in random.sample(range(1, 28), 5):
        c.execute('''INSERT INTO menu (food, day, date) VALUES ({:d}, '{:s}', date('now'))'''.format(f, d))


conn.commit()

conn.close()
