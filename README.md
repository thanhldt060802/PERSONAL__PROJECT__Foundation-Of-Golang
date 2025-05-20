<h1>Các thí nghiệm với Ergo trong xử lý đồng thời theo Actor model</h1>
<h3>Ví dụ 1:</h3>
<span>Tình huống 2 Sender call đồng bộ đến local Receiver trên cùng node1 và call đến remote Receiver trên node2.</span>
<ul>
<li>Không đóng gói hệ thống.</li>
<li>Sender được cung cấp tên của Receiver.</li>
<li>2 Sender đồng thời call đến lần lượt local Receiver và remote Receiver bởi việc Sender tự trigger chính nó để thay đổi vị trí (local/remote) Receiver cần call.</li>
</ul>
<br>                                      
<h3>Ví dụ 2:</h3>
<span>Tình huống n Sender call đồng bộ đến n local Receiver trên cùng node1 và call đến n remote Receiver trên node2 tương ứng.</span>
<ul>
<li>Đóng gói hệ thống.</li>
<li>Sender x được cung cấp tên của Receiver x để chỉ tương tác với Receiver x.</li>
<li>n Sender đồng thời call đến lần lượt n local Receiver và n remote Receiver bởi việc Sender x tự trigger chính nó để thay đổi vị trí (local/remote) Receiver x cần call.</li>
</ul>
<br>
<h3>Ví dụ 3:</h3>
<span>Tình huống 1 Sender trên node1 call bất đồng bộ đến n remote Receiver trên node2.</span>
<ul>
<li>Đóng gói hệ thống.</li>
<li>n Receiver được cung cấp tên của Sender để phát tín hiệu "idle" cho Sender, Sender nhận đúng tính hiệu và send qua cho Receiver gửi tín hiệu.</li>
<li>n Receiver đồng thời send đến Sender bởi khi Receiver được khởi tạo hoặc Receiver hoàn tất xử lý.</li>
<li>n Receiver được quản lý bởi 1 Supervisor để restart lại Receiver khi nó bị crash bởi tình huống crash giả lập.</li>
</ul>
<br>
<h3>Ví dụ 4:</h3>
<span>Tình huống vận dụng: 1 Sender trên node1 call bất đồng bộ đến n remote Receiver trên node2, mỗi Receiver xử lý 1 task và thực hiện cập nhật lên database khi hoàn tất hoặc gặp lỗi.</span>
<ul>
<li>Đóng gói hệ thống.</li>
<li>n Receiver được cung cấp tên của Sender để phát tín hiệu "idle" cho Sender, Sender nhận đúng tính hiệu và send task qua cho Receiver gửi tín hiệu.</li>
<li>n Receiver đồng thời send đến Sender bởi khi Receiver được khởi tạo hoặc Receiver hoàn tất xử lý task.</li>
<li>Dữ liệu task mà Receiver xử lý sẽ được cập nhật lại ngay khi Receiver hoàn tất xử lý hoặc gặp lỗi trong quá trình xử lý.</li>
<li>n Receiver được quản lý bởi 1 Supervisor để restart lại Receiver khi nó bị crash bởi tình huống crash giả lập.</li>
</ul>
<br>
<h3>Ví dụ 5:</h3>
<span>Triển khai FSM theo mô hình có sẵn của Erlang.</span>
<ul>
<li>Không đóng gói hệ thống.</li>
<li>Định nghĩa các state và quy tắc chuyển state theo event.</li>
</ul>