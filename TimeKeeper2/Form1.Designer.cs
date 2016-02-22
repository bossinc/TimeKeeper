namespace TimeKeeper2
{
    partial class MainForm
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(MainForm));
            this.timer = new System.Windows.Forms.Timer(this.components);
            this.tvActiveWindows = new System.Windows.Forms.TreeView();
            this.menuStrip = new System.Windows.Forms.MenuStrip();
            this.fileToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.newToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.openToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.editToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.deleteSelectionToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.tbNotes = new System.Windows.Forms.TextBox();
            this.lblNotes = new System.Windows.Forms.Label();
            this.lblEndingTime = new System.Windows.Forms.Label();
            this.label5 = new System.Windows.Forms.Label();
            this.lblStartingTime = new System.Windows.Forms.Label();
            this.label6 = new System.Windows.Forms.Label();
            this.label7 = new System.Windows.Forms.Label();
            this.lblTime = new System.Windows.Forms.Label();
            this.lblProjectName = new System.Windows.Forms.Label();
            this.dgvEntries = new System.Windows.Forms.DataGridView();
            this.Column1 = new System.Windows.Forms.DataGridViewTextBoxColumn();
            this.btnReview = new System.Windows.Forms.Button();
            this.btnStart = new System.Windows.Forms.Button();
            this.cbDrawing = new System.Windows.Forms.CheckBox();
            this.menuStrip.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.dgvEntries)).BeginInit();
            this.SuspendLayout();
            // 
            // timer
            // 
            this.timer.Tick += new System.EventHandler(this.timer_Tick);
            // 
            // tvActiveWindows
            // 
            this.tvActiveWindows.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.tvActiveWindows.Location = new System.Drawing.Point(207, 84);
            this.tvActiveWindows.Name = "tvActiveWindows";
            this.tvActiveWindows.Size = new System.Drawing.Size(565, 181);
            this.tvActiveWindows.TabIndex = 3;
            // 
            // menuStrip
            // 
            this.menuStrip.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.fileToolStripMenuItem,
            this.editToolStripMenuItem});
            this.menuStrip.Location = new System.Drawing.Point(0, 0);
            this.menuStrip.Name = "menuStrip";
            this.menuStrip.Size = new System.Drawing.Size(784, 24);
            this.menuStrip.TabIndex = 4;
            this.menuStrip.Text = "menuStrip";
            // 
            // fileToolStripMenuItem
            // 
            this.fileToolStripMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.newToolStripMenuItem,
            this.openToolStripMenuItem});
            this.fileToolStripMenuItem.Name = "fileToolStripMenuItem";
            this.fileToolStripMenuItem.Size = new System.Drawing.Size(37, 20);
            this.fileToolStripMenuItem.Text = "File";
            // 
            // newToolStripMenuItem
            // 
            this.newToolStripMenuItem.Name = "newToolStripMenuItem";
            this.newToolStripMenuItem.Size = new System.Drawing.Size(103, 22);
            this.newToolStripMenuItem.Text = "New";
            this.newToolStripMenuItem.Click += new System.EventHandler(this.newToolStripMenuItem_Click);
            // 
            // openToolStripMenuItem
            // 
            this.openToolStripMenuItem.Name = "openToolStripMenuItem";
            this.openToolStripMenuItem.Size = new System.Drawing.Size(103, 22);
            this.openToolStripMenuItem.Text = "Open";
            this.openToolStripMenuItem.Click += new System.EventHandler(this.openToolStripMenuItem_Click);
            // 
            // editToolStripMenuItem
            // 
            this.editToolStripMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.deleteSelectionToolStripMenuItem});
            this.editToolStripMenuItem.Name = "editToolStripMenuItem";
            this.editToolStripMenuItem.Size = new System.Drawing.Size(39, 20);
            this.editToolStripMenuItem.Text = "Edit";
            // 
            // deleteSelectionToolStripMenuItem
            // 
            this.deleteSelectionToolStripMenuItem.Name = "deleteSelectionToolStripMenuItem";
            this.deleteSelectionToolStripMenuItem.Size = new System.Drawing.Size(158, 22);
            this.deleteSelectionToolStripMenuItem.Text = "Delete Selection";
            this.deleteSelectionToolStripMenuItem.Click += new System.EventHandler(this.deleteSelectionToolStripMenuItem_Click);
            // 
            // tbNotes
            // 
            this.tbNotes.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.tbNotes.Location = new System.Drawing.Point(207, 284);
            this.tbNotes.Multiline = true;
            this.tbNotes.Name = "tbNotes";
            this.tbNotes.ScrollBars = System.Windows.Forms.ScrollBars.Both;
            this.tbNotes.Size = new System.Drawing.Size(565, 145);
            this.tbNotes.TabIndex = 5;
            // 
            // lblNotes
            // 
            this.lblNotes.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.lblNotes.AutoSize = true;
            this.lblNotes.Location = new System.Drawing.Point(207, 268);
            this.lblNotes.Name = "lblNotes";
            this.lblNotes.Size = new System.Drawing.Size(38, 13);
            this.lblNotes.TabIndex = 6;
            this.lblNotes.Text = "Notes:";
            // 
            // lblEndingTime
            // 
            this.lblEndingTime.AutoSize = true;
            this.lblEndingTime.Location = new System.Drawing.Point(657, 55);
            this.lblEndingTime.Name = "lblEndingTime";
            this.lblEndingTime.Size = new System.Drawing.Size(0, 13);
            this.lblEndingTime.TabIndex = 19;
            // 
            // label5
            // 
            this.label5.AutoSize = true;
            this.label5.Location = new System.Drawing.Point(579, 55);
            this.label5.Name = "label5";
            this.label5.Size = new System.Drawing.Size(55, 13);
            this.label5.TabIndex = 18;
            this.label5.Text = "End Time:";
            // 
            // lblStartingTime
            // 
            this.lblStartingTime.AutoSize = true;
            this.lblStartingTime.Location = new System.Drawing.Point(657, 33);
            this.lblStartingTime.Name = "lblStartingTime";
            this.lblStartingTime.Size = new System.Drawing.Size(0, 13);
            this.lblStartingTime.TabIndex = 17;
            // 
            // label6
            // 
            this.label6.AutoSize = true;
            this.label6.Location = new System.Drawing.Point(579, 33);
            this.label6.Name = "label6";
            this.label6.Size = new System.Drawing.Size(58, 13);
            this.label6.TabIndex = 16;
            this.label6.Text = "Start Time:";
            // 
            // label7
            // 
            this.label7.AutoSize = true;
            this.label7.Location = new System.Drawing.Point(415, 33);
            this.label7.Name = "label7";
            this.label7.Size = new System.Drawing.Size(33, 13);
            this.label7.TabIndex = 15;
            this.label7.Text = "Time:";
            // 
            // lblTime
            // 
            this.lblTime.AutoSize = true;
            this.lblTime.Font = new System.Drawing.Font("Microsoft Sans Serif", 27.75F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.lblTime.Location = new System.Drawing.Point(411, 39);
            this.lblTime.Name = "lblTime";
            this.lblTime.Size = new System.Drawing.Size(143, 42);
            this.lblTime.TabIndex = 14;
            this.lblTime.Text = "0:00:00";
            // 
            // lblProjectName
            // 
            this.lblProjectName.AutoSize = true;
            this.lblProjectName.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.lblProjectName.Location = new System.Drawing.Point(212, 33);
            this.lblProjectName.Name = "lblProjectName";
            this.lblProjectName.Size = new System.Drawing.Size(0, 20);
            this.lblProjectName.TabIndex = 22;
            // 
            // dgvEntries
            // 
            this.dgvEntries.AllowUserToAddRows = false;
            this.dgvEntries.AllowUserToDeleteRows = false;
            this.dgvEntries.AllowUserToResizeColumns = false;
            this.dgvEntries.AllowUserToResizeRows = false;
            this.dgvEntries.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left)));
            this.dgvEntries.AutoSizeColumnsMode = System.Windows.Forms.DataGridViewAutoSizeColumnsMode.Fill;
            this.dgvEntries.BackgroundColor = System.Drawing.SystemColors.ControlLight;
            this.dgvEntries.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.AutoSize;
            this.dgvEntries.ColumnHeadersVisible = false;
            this.dgvEntries.Columns.AddRange(new System.Windows.Forms.DataGridViewColumn[] {
            this.Column1});
            this.dgvEntries.Location = new System.Drawing.Point(12, 27);
            this.dgvEntries.Name = "dgvEntries";
            this.dgvEntries.ReadOnly = true;
            this.dgvEntries.RowHeadersVisible = false;
            this.dgvEntries.RowHeadersWidthSizeMode = System.Windows.Forms.DataGridViewRowHeadersWidthSizeMode.DisableResizing;
            this.dgvEntries.Size = new System.Drawing.Size(189, 373);
            this.dgvEntries.TabIndex = 23;
            this.dgvEntries.TabStop = false;
            // 
            // Column1
            // 
            this.Column1.HeaderText = "Column1";
            this.Column1.Name = "Column1";
            this.Column1.ReadOnly = true;
            // 
            // btnReview
            // 
            this.btnReview.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.btnReview.Location = new System.Drawing.Point(12, 406);
            this.btnReview.Name = "btnReview";
            this.btnReview.Size = new System.Drawing.Size(189, 23);
            this.btnReview.TabIndex = 24;
            this.btnReview.Text = "Review Selection";
            this.btnReview.UseVisualStyleBackColor = true;
            this.btnReview.Click += new System.EventHandler(this.btnReview_Click);
            // 
            // btnStart
            // 
            this.btnStart.Location = new System.Drawing.Point(207, 58);
            this.btnStart.Name = "btnStart";
            this.btnStart.Size = new System.Drawing.Size(75, 23);
            this.btnStart.TabIndex = 25;
            this.btnStart.Text = "Start";
            this.btnStart.UseVisualStyleBackColor = true;
            this.btnStart.Click += new System.EventHandler(this.btnStart_Click);
            // 
            // cbDrawing
            // 
            this.cbDrawing.AutoSize = true;
            this.cbDrawing.Location = new System.Drawing.Point(289, 61);
            this.cbDrawing.Name = "cbDrawing";
            this.cbDrawing.Size = new System.Drawing.Size(65, 17);
            this.cbDrawing.TabIndex = 26;
            this.cbDrawing.Text = "Drawing";
            this.cbDrawing.UseVisualStyleBackColor = true;
            // 
            // MainForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(784, 441);
            this.Controls.Add(this.cbDrawing);
            this.Controls.Add(this.btnStart);
            this.Controls.Add(this.btnReview);
            this.Controls.Add(this.dgvEntries);
            this.Controls.Add(this.lblProjectName);
            this.Controls.Add(this.lblEndingTime);
            this.Controls.Add(this.label5);
            this.Controls.Add(this.lblStartingTime);
            this.Controls.Add(this.label6);
            this.Controls.Add(this.label7);
            this.Controls.Add(this.lblTime);
            this.Controls.Add(this.lblNotes);
            this.Controls.Add(this.tbNotes);
            this.Controls.Add(this.tvActiveWindows);
            this.Controls.Add(this.menuStrip);
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.MainMenuStrip = this.menuStrip;
            this.MinimumSize = new System.Drawing.Size(800, 360);
            this.Name = "MainForm";
            this.Text = "Time Keeper 2";
            this.menuStrip.ResumeLayout(false);
            this.menuStrip.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.dgvEntries)).EndInit();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Timer timer;
        private System.Windows.Forms.TreeView tvActiveWindows;
        private System.Windows.Forms.MenuStrip menuStrip;
        private System.Windows.Forms.ToolStripMenuItem fileToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem newToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem openToolStripMenuItem;
        private System.Windows.Forms.TextBox tbNotes;
        private System.Windows.Forms.Label lblNotes;
        private System.Windows.Forms.Label lblEndingTime;
        private System.Windows.Forms.Label label5;
        private System.Windows.Forms.Label lblStartingTime;
        private System.Windows.Forms.Label label6;
        private System.Windows.Forms.Label label7;
        private System.Windows.Forms.Label lblTime;
        private System.Windows.Forms.Label lblProjectName;
        private System.Windows.Forms.DataGridView dgvEntries;
        private System.Windows.Forms.DataGridViewTextBoxColumn Column1;
        private System.Windows.Forms.Button btnReview;
        private System.Windows.Forms.Button btnStart;
        private System.Windows.Forms.ToolStripMenuItem editToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem deleteSelectionToolStripMenuItem;
        private System.Windows.Forms.CheckBox cbDrawing;
    }
}

