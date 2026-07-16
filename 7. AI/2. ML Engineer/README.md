# ML Engineer yo'li

Go backend fonidan ML Engineer kasbiga o'tish uchun o'quv reja. Repodagi mavjud bilimlar (Database, Algorithm, Linux, DevOps, System Design) bu yo'lning infratuzilma qismini yopadi — bu papka esa yetishmayotgan ML qismini yopadi.

Tezkor ish topish yo'li — `../1. AI Engineer/` (LLM integratsiya, 2–4 oy). Bu papka esa uzoq muddatli chuqurlik: model qurish, o'rgatish va production'da yuritish. Tavsiya etilgan ketma-ketlik `../README.md` da.

## Struktura va tavsiya etilgan tartib

```
1. Python ecosystem/     ← boshlanish nuqtasi (ikkala yo'lga poydevor)
   1. Python             — til asoslari (Basics + Advanced tayyor)
   2. NumPy              — vektorlashgan hisob-kitob
   3. Pandas             — data bilan ishlash (SQL bilimiga bog'lanadi)
   4. Matplotlib         — visualization
   5. Scikit-learn       — klassik ML kutubxonasi
   6. PyTorch            — deep learning framework

2. Matematika/           ← Python bilan parallel yursa bo'ladi
   0. Poydevor           — maktab matematikasi noldan (Khan Academy, reja ichida)
   1. Linear algebra     — vektor, matritsa, matrix multiplication
   2. Calculus           — hosila, gradient (gradient descent uchun)
   3. Probability        — ehtimollik, taqsimotlar
   4. Statistics         — hypothesis testing, evaluation asoslari

3. Machine learning/     ← 1 va 2 dan keyin
   1. Fundamentals       — train/test, overfitting, metrics, feature engineering
   2. Supervised learning
   3. Unsupervised learning
   4. Deep learning
   5. MLOps              — model deploy, monitoring, drift (DevOps bilimiga bog'lanadi)

4. Data Engineering/     ← ML fundamentals'dan keyin darhol kerak
   1. Data pipelines (Airflow)  — orchestration (Temporal bilimiga bog'lanadi)
   2. Spark              — katta hajmdagi data'ni qayta ishlash
   3. Kafka streaming    — real-time data oqimi
   4. Data formats (Parquet, Arrow) — column-oriented formatlar
   5. Feature store      — feature'larni saqlash va qayta ishlatish

5. LLM internals/        ← Deep learning'dan keyin: model ichini tushunish
   1. Transformers       — zamonaviy modellarning arxitekturasi
   2. Fine-tuning        — tayyor modelni moslashtirish
   (amaliy qismi — Embeddings, Vector DB, RAG, Agents — `../1. AI Engineer/` da)

6. ML System Design/     ← interview bosqichi (System Design bilimiga bog'lanadi)
   Recommender, search, fraud detection kabi tizimlar dizayni

7. Cloud/                ← deploy bosqichida kerak; training o'rganayotganda Colab/local GPU yetadi
   1. AWS
   2. GCP
   3. Kubernetes         — production ML deyarli har doim K8s ustida

8. Security/             ← yakuniy bosqich — model va LLM bilimini talab qiladi
   1. Adversarial attacks
   2. Prompt injection

9. Database/             ← poydevor (repo root'dan ko'chirilgan, o'rganib bo'lingan)
   Postgres (Basic + Advanced internals), Redis — SQL bu kasbning №1 kundalik tool'i

10. Algorithm/           ← poydevor (repo root'dan ko'chirilgan)
    Nazariya (18 mavzu) + 100 kunlik LeetCode challenge
```

Papka raqamlari = o'rganish prioriteti (1–8). 9–10 esa tayyor poydevor, tartibda qatnashmaydi.

## Amaliyot

Alohida "Amaliyot" papkasi yo'q — har bo'limning amaliyoti o'z ichida bo'ladi (masalan, Machine learning ichida train loyihalari, Data Engineering ichida pipeline loyihalari). Maqsad: nazariya va mashq yonma-yon tursin.

## Repodagi boshqa umumiy poydevor

Linux, DevOps, System Design, Golang, Network — backend va ML yo'liga birdek xizmat qiladi, repo root'da qoladi. Database va Algorithm esa 2026-07-12 da shu papkaga ko'chirilgan (9 va 10).

## Til va format

- Tushuntirishlar o'zbek tilida, texnik atamalar ingliz tilida
- Diagrammalar — Mermaid formatda
