
 **Step 1: TCP Listener ishga tushirish**
 
- Port: `****`
- Har bir yangi TCP ulanish uchun `Session` ochiladi
- func: listenSMPP


**Step 2: SMPP sessiyasini yuritish** uchun asosiy handler
- func: handleSMPPConnection
	- func: NewSession
	- ReadPDU
		- sendEnquireLink
		- sendUnbind
	- handlePDU


 Step 2: `bind_transceiver` ni qabul qilish

- `system_id`, `password`, `addr_ton`, `addr_npi` validatsiya qilinadi
- Javob sifatida `bind_transceiver_resp` qaytaring

### 📌 Step 3: `submit_sm` ni qabul qilish

- Logga yozilsin (`from`, `to`, `text`)
    
- Javob: `submit_sm_resp` (`message_id` bilan)
    
- 1–2 sekunddan so‘ng `deliver_sm` yuborilsin
    

### 📌 Step 4: `deliver_sm` yuborish

- Qabul qilingan `submit_sm` ni “clientga qaytarish”
    
- Client `deliver_sm_resp` bilan javob beradi
    

### 📌 Step 5: `enquire_link` va `unbind` qo‘llab-quvvatlash