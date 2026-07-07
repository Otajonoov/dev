

![](../../assets/obsidian-images/Pasted%20image%2020250714155736.png)
#### Metrikalarning turlari:

1.  **==Counter== (Hisoblagich)**: Faqat oshib boradigan qiymat. Reset bo'lishi faqat restart bilan
    - HTTP so'rovlar soni: `http_requests_total`
    - Xatoliklar soni: `http_errors_total`
    - Kafka xabarlar soni: `kafka_messages_produced_total`

2. **==Gauge== (Ko‘rsatkich)**: Ortishi ham, kamayishi ham mumkin
    - CPU:  `cpu_usage_percent`
    - RAM: `memory_usage_bytes`
    - Connectionlar soni: `active_connections`

3. **==Histogram==**: Qiymatlarni **diapazonlarga** (buckets) ajratib, ularning taqsimotini ko‘rsatadi.
- **Qo‘llaniladi**: Resurs ishlash vaqtlarini (latency) o‘lchashda.

    - HTTP javob vaqtlar: `http_request_duration_seconds`
    - Tizim ishga tushish vaqti
    
- **Qo‘shimcha qiymatlar**:
    
    - `_count`: nechta so‘rov bo‘lgan
        
    - `_sum`: umumiy vaqt
        
    - `_bucket`: tanlangan vaqt oralig‘ida nechta so‘rov bo‘lgan

4. **==Summary==**: Histogramga o‘xshash, lekin **quantile** (foizlar) asosida aniqlik beradi.

    - 95% so‘rovlar 500ms dan kam davom etdi
        
- **Kamchiligi**:
    - Kengaytirilgan statistikani qo‘llab-quvvatlamaydi, shuning uchun Prometheusda ko‘p ishlatilmaydi

5. Custom/Derived Metrics (Hosila metrikalar): - Bular mavjud metrikalardan chiqarilgan qiymatlar (foydalanish foizi, o‘rtacha vaqt, ratio).
```promql
rate(http_errors_total[5m]) / rate(http_requests_total[5m]) 
```
- → 5 daqiqadagi xatolik foizi

| Metrika turi          | Tavsif                                      |
| --------------------- | ------------------------------------------- |
| **Performance**       | Tezlik, throughput, latency, error rate     |
| **Resource usage**    | CPU, RAM, Disk, Network                     |
| **Application-level** | API chaqiriqlari, login muvaffaqiyati, TPS  |
| **Business metrics**  | Foydalanuvchilar soni, sotuvlar, konversiya |
| **Infrastructure**    | Tarmoq yuklamasi, container holati, podlar  |






[Prometheus](https://prometheus.io/) для работы с метриками. Он включает в себя:
- сервер — хранилка и сборщик метрик;
- формат данных;
- язык запросов — еще его называют PromQL.


