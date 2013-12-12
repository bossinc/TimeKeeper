using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Windows.Forms;
using System.IO;
using System.Runtime.InteropServices;

namespace TimeKeeper2
{
    public partial class MainForm : Form
    {
        [DllImport("user32.dll")]
        static extern int GetForegroundWindow();

        [DllImport("user32.dll")]
        static extern int GetWindowText(int hWnd, StringBuilder text, int count);

        [DllImport("user32.dll")]
        static extern bool GetLastInputInfo(out LASTINPUTINFO plii);

        [StructLayout(LayoutKind.Sequential)]
        struct LASTINPUTINFO
        {
            public static readonly int SizeOf =
                   Marshal.SizeOf(typeof(LASTINPUTINFO));

            [MarshalAs(UnmanagedType.U4)]
            public int cbSize;
            [MarshalAs(UnmanagedType.U4)]
            public int dwTime;
        }

        private string captionWindowLabel;
        private string idWindowLabel;

        private List<string> activeWindowLabels;
        private List<int> activeWindowTimes;
        private int lastIndexTicked;

        private static int AFKTimeAmout = 300000;
        private bool dnuisAFK;
        private bool isAFK
        {
            get
            {
                return dnuisAFK;
            }
            set
            {
                if (!dnuisAFK && value)
                {
                    activeWindowTimes[lastIndexTicked] -= AFKTimeAmout;
                    TotalTimeWorking -= AFKTimeAmout;
                    btnStart_Click(null, EventArgs.Empty);
                    MessageBox.Show("You are AFK. The past 5 minutes have been removed as time spent working, the counter has stopped and the file has been saved.",
                        "Warning",
                        MessageBoxButtons.OK,
                        MessageBoxIcon.Warning,
                        MessageBoxDefaultButton.Button1);
                }
                dnuisAFK = value;
            }
        }

        private int TotalTimeWorking;

        private string SaveFileText;
        private List<string> SaveEntriesText;
        private string fileLocation;

        private int lastUpdateValueTV;

        public MainForm()
        {
            SaveFileText = "";
            fileLocation = "";
            InitializeComponent();
        }

        protected override void OnClosing(CancelEventArgs e)
        {
            if(fileLocation != "")
                SaveFile();
            base.OnClosing(e);
        }

        private void InitializeVaribles()
        {
            captionWindowLabel = "";
            idWindowLabel = "";
            activeWindowLabels = new List<string>();
            activeWindowTimes = new List<int>();
            lastIndexTicked = -1;
            tvActiveWindows.Nodes.Clear();
            tbNotes.Text = "";
            TotalTimeWorking = 0;
            lastUpdateValueTV = 0;
            isAFK = false;
            lblStartingTime.Text = "";
            lblTime.Text = "0:00:00";
        }

        private void NewFile()
        {
            InitializeVaribles();
            dgvEntries.Rows.Clear();
            SaveEntriesText = new List<string>();
            activeWindowLabels.Add("Dead Ticks");
            activeWindowTimes.Add(0);
            lblStartingTime.Text = "";
            lblEndingTime.Text = "";
        }

        private void OpenNewFile()
        {
            OpenFileDialog fileDialog = new OpenFileDialog();

            // Default to the directory which contains our content files.
            string relativePath = Path.Combine("../../../../Content");
            string contentPath = Path.GetFullPath(relativePath);

            fileDialog.InitialDirectory = contentPath;

            fileDialog.Title = "Open - Time Keeper 2";

            fileDialog.Filter = "Time Keeper Files 2 (*.tk2)|*.tk2|" + "All Files (*.*)|*.*";

            if (fileDialog.ShowDialog() == DialogResult.OK)
            {
                fileLocation = fileDialog.FileName;
                this.Text = removeFilePath(fileLocation) + " - Time Keeper 2";
                OpenFile();
            }
            else
                return;

        }

        private void OpenFile()
        {
            List<byte> tempList = new List<byte>();
            try
            {
                using (System.IO.StreamReader file = new System.IO.StreamReader(fileLocation))
                {
                    string temp = file.ReadToEnd();
                    for (int i = 0; temp.Length > i; i++)
                    {
                        tempList.Add((byte)temp[i]);
                    }
                }
            }
            catch (System.ArgumentNullException)
            {
                MessageBox.Show("File not found.", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                return;
            }
            catch (System.ArgumentException)
            {
                return;
            }

            byte[] encSave = new byte[tempList.Count];
            for (int i = 0; tempList.Count > i; i++)
            {
                encSave[i] = tempList[i];
            }
            //decrypts file
            try
            {
                EnigmaForce ef = new EnigmaForce();
                byte[] encPass = (byte[])encSave;
                string rawPass = ef.Decrypt(encPass);
                SaveFileText = rawPass;
            }
            catch (Exception)
            {
                MessageBox.Show("File is corrupted.", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error, MessageBoxDefaultButton.Button1);
                newToolStripMenuItem_Click(null, EventArgs.Empty);
                return;
            }

            NewFile();
            LoadFile();
            lblProjectName.Text = removeFilePath(fileLocation);
            lblStartingTime.Text = DateTime.Now.ToString();
        }

        public void SaveFile()
        {
            timer.Stop();
            lblEndingTime.Text = DateTime.Now.ToString();

            SaveFileText += lblStartingTime.Text + "\r\n";
            SaveFileText += lblEndingTime.Text + "\r\n";
            SaveFileText += tbNotes.Text + "\r\n-SN%^1-\r\n";
            for(int i = 0; activeWindowLabels.Count > i; i++)
            {
                SaveFileText += activeWindowLabels[i] + "\r\n";
                SaveFileText += activeWindowTimes[i] + "\r\n";
            }
            SaveFileText += "-End%^1-";

            //Encrypts SaveFileText
            EnigmaForce ef = new EnigmaForce();
            byte[] encSave = ef.Encrypt(SaveFileText);

            if (fileLocation == "")
            {
                SaveFileDialog fileDialog = new SaveFileDialog();

                string relativePath = Path.Combine("../../../../Content");
                string contentPath = Path.GetFullPath(relativePath);

                fileDialog.InitialDirectory = contentPath;

                fileDialog.Title = "Save - Time Keeper 2";

                fileDialog.Filter = "Time Keeper 2 Files (*.tk2)|*.tk2";


                if (fileDialog.ShowDialog() == DialogResult.OK)
                {
                    fileLocation = fileDialog.FileName;
                    lblProjectName.Text = removeFilePath(fileLocation);
                }
                else
                {
                    return;
                }
            }

            //Writes SaveFileText to fileLocation
            using (System.IO.StreamWriter file = new System.IO.StreamWriter(fileLocation))
            {
                string temp = "";
                foreach (byte b in encSave)
                {
                    temp += (char)b;
                }
                file.Write(temp);
            }
        }

        private void LoadFile()
        {
            string[] splits = SaveFileText.Split(new string[] { "-End%^1-" }, StringSplitOptions.RemoveEmptyEntries);
            foreach (string s in splits)
            {
                SaveEntriesText.Add(s);
                dgvEntries.Rows.Add(s/*.Split(new string[] { "\r\n" }, StringSplitOptions.None)[0]*/);
            }
        }

        private string removeFilePath(string filePath)
        {
            //Removes file path and displays the name of the file
            int indexlength = filePath.Length;
            for (int i = 0; indexlength > i; i++)
            {
                if (filePath[i] == '\\')
                {
                    filePath = filePath.Remove(0, i);
                    indexlength -= i;
                    i = 0;
                }
            }
            filePath = filePath.Remove(0, 1);
            filePath = filePath.Remove(filePath.Length - 4, 4);

            return filePath;
        }

        private void GetActiveWindow()
        {
            const int nChars = 256;
            int handle = 0;
            StringBuilder Buff = new StringBuilder(nChars);

            handle = GetForegroundWindow();

            if (GetWindowText(handle, Buff, nChars) > 0)
            {
                this.captionWindowLabel = Buff.ToString();
                this.idWindowLabel = handle.ToString();
                for (int s = 0; s < activeWindowLabels.Count; s++)
                {
                    if (activeWindowLabels[s] == captionWindowLabel)
                    {
                        activeWindowTimes[s] += timer.Interval;
                        lastIndexTicked = s;
                        return;
                    }
                }
                if (captionWindowLabel.Contains("porn") || captionWindowLabel.Contains("Porn") || captionWindowLabel.Contains("PORN") || captionWindowLabel.Contains("nsfw"))
                {
                    if (fileLocation != "")
                        SaveFile();
                    else
                    {
                        timer.Stop();
                    }
                    MessageBox.Show("Does not keep track of porn.", "Warning", MessageBoxButtons.OK, MessageBoxIcon.Warning, MessageBoxDefaultButton.Button1);
                    return;
                }
                //makes sure it is not a loading form
                if (captionWindowLabel.Contains("% complete"))
                {
                    return;
                }
                activeWindowLabels.Add(captionWindowLabel);
                activeWindowTimes.Add(timer.Interval);
            }
            else
            {
                //makes sure it is not a loading form
                if (captionWindowLabel.Contains("% com"))
                {
                    activeWindowTimes[0] += timer.Interval;
                }
            }

        }

        private void ShowActiveWindowsInTreeVeiw()
        {
            for (int i = lastUpdateValueTV; i < activeWindowLabels.Count; i++)
            {
                string[] labelSplits = activeWindowLabels[i].Split(new string[] { " - " },  StringSplitOptions.None);

                //Creates base nodes
                int hiIndex = CreateBaseNodes(labelSplits);

                //fills base nodes
                FillNodes(labelSplits, hiIndex);
                lastUpdateValueTV = i + 1;

            }

            //update time
            List<List<int>> rawnodeIndex2 = new List<List<int>>();
            List<int> rawnodeIndex1 = new List<int>();
            for (int i = 0; i < activeWindowLabels.Count; i++)
            {
                List<int> nodeIndex = GetNode(activeWindowLabels[i]);
                
                if(nodeIndex.Count == 3)
                {
                    tvActiveWindows.Nodes[nodeIndex[0]].Nodes[nodeIndex[1]].Nodes[nodeIndex[2]].Text = tvActiveWindows.Nodes[nodeIndex[0]].Nodes[nodeIndex[1]].Nodes[nodeIndex[2]].Text.Split(new string[] { " : " }, StringSplitOptions.None)[0] + " : " + FormatTime(activeWindowTimes[i]);
                }
                if(nodeIndex.Count == 2)
                {
                    string[] tempSubSplits = tvActiveWindows.Nodes[nodeIndex[0]].Nodes[nodeIndex[1]].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                    string tempName = tempSubSplits[0];
                    for (int i3 = 1; tempSubSplits.Length - 1 > i3; i3++)
                    {
                        tempName += " : " + tempSubSplits[i3];
                    }
                    tvActiveWindows.Nodes[nodeIndex[0]].Nodes[nodeIndex[1]].Text = tempName + " : " + FormatTime(activeWindowTimes[i]);
                    if (tvActiveWindows.Nodes[nodeIndex[0]].Nodes[nodeIndex[1]].Nodes.Count > 0)
                    {
                        rawnodeIndex2.Add(nodeIndex);
                    }

                }
                if(nodeIndex.Count == 1)
                {
                    tvActiveWindows.Nodes[nodeIndex[0]].Text = tvActiveWindows.Nodes[nodeIndex[0]].Text.Split(new string[] { " : " }, StringSplitOptions.None)[0] + " : " + FormatTime(activeWindowTimes[i]);
                    if (tvActiveWindows.Nodes[nodeIndex[0]].Nodes.Count > 0)
                    {
                        rawnodeIndex1.Add(nodeIndex[0]);
                    }
                }
            }

            for (int i = 0; i < tvActiveWindows.Nodes.Count; i++)
            {
                //foreach (int si in rawnodeIndex1)
                //{
                //    if (si == i)
                //    {
                //        timeTotal += DeformatTime(tvActiveWindows.Nodes[i].Text.Split(new string[] { " : " }, StringSplitOptions.None)[1]);
                //        break;
                //    }
                //}
                for (int i2 = 0; i2 < tvActiveWindows.Nodes[i].Nodes.Count; i2++)
                {
                    int subTimeTotal = 0;
                    string[] tempSubSplits = tvActiveWindows.Nodes[i].Nodes[i2].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                    foreach (List<int> si in rawnodeIndex2)
                    {
                        if (si[0] == i && si[1] == i2)
                        {
                            subTimeTotal += DeformatTime(tempSubSplits[tempSubSplits.Length - 1]);
                            break;
                        }
                    }

                    for (int i3 = 0; i3 < tvActiveWindows.Nodes[i].Nodes[i2].Nodes.Count; i3++)
                    {
                        tempSubSplits = tvActiveWindows.Nodes[i].Nodes[i2].Nodes[i3].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                        subTimeTotal += DeformatTime(tempSubSplits[tempSubSplits.Length - 1]);
                    }

                    if (subTimeTotal > 0)
                    {
                        tempSubSplits = tvActiveWindows.Nodes[i].Nodes[i2].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                        string tempName = tempSubSplits[0];
                        for (int i3 = 1; tempSubSplits.Length - 1 > i3; i3++)
                        {
                            tempName += " : " + tempSubSplits[i3];
                        }
                        tvActiveWindows.Nodes[i].Nodes[i2].Text = tempName + " : " + FormatTime(subTimeTotal);
                    }
                }

                //BaseNodes
                int timeTotal = 0;
                string[] tempSplits = tvActiveWindows.Nodes[i].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                foreach (int si in rawnodeIndex1)
                {
                    if (si == i)
                    {
                        timeTotal += DeformatTime(tempSplits[tempSplits.Length - 1]);
                        break;
                    }
                }
                for (int i2 = 0; i2 < tvActiveWindows.Nodes[i].Nodes.Count; i2++)
                {
                    tempSplits = tvActiveWindows.Nodes[i].Nodes[i2].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                    timeTotal += DeformatTime(tempSplits[tempSplits.Length - 1]);
                    
                }
                if (timeTotal > 0)
                {
                    tempSplits = tvActiveWindows.Nodes[i].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                    string tempName = tempSplits[0];
                    for (int i3 = 1; tempSplits.Length - 1 > i3; i3++)
                    {
                        tempName += " : " + tempSplits[i3];
                    }
                    tvActiveWindows.Nodes[i].Text = tempName + " : " + FormatTime(timeTotal);
                    
                }
            }
        }

        private List<int> GetNode(string name)
        {
            List<int> rList = new List<int>();
            string[] sName = name.Split(new string[] { " - " }, StringSplitOptions.None);
            int nodeLayer = sName.Length-1;
            if (nodeLayer > 2)
            {
                for (int i = nodeLayer-2; i > 0; i--)
                {
                    sName[0] += " - " + sName[i];
                }
                sName[nodeLayer - 2] = sName[0];
            }

            for (int i = 0; i < tvActiveWindows.Nodes.Count; i++)
            {
                string[] tempSplits = tvActiveWindows.Nodes[i].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                string tempName = tempSplits[0];
                for (int i3 = 1; tempSplits.Length - 1 > i3; i3++)
                {
                    tempName += " : " + tempSplits[i3];
                }
                if (sName[sName.Length-1] == tempName)
                {
                    rList.Add(i);

                    if (nodeLayer >= 1)
                    {
                        for (int i2 = 0; i2 < tvActiveWindows.Nodes[i].Nodes.Count; i2++)
                        {
                            tempSplits = tvActiveWindows.Nodes[i].Nodes[i2].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                            tempName = tempSplits[0];
                            for (int i3 = 1; tempSplits.Length - 1 > i3; i3++)
                            {
                                tempName += " : " + tempSplits[i3];
                            }
                            if (sName[sName.Length - 2] == tempName)
                            {
                                rList.Add(i2);

                                if (nodeLayer >= 2)
                                {
                                    for (int i3 = 0; i3 < tvActiveWindows.Nodes[i].Nodes[i2].Nodes.Count; i3++)
                                    {
                                        tempSplits = tvActiveWindows.Nodes[i].Nodes[i2].Nodes[i3].Text.Split(new string[] { " : " }, StringSplitOptions.None);
                                        tempName = tempSplits[0];
                                        for (int i4 = 1; tempSplits.Length - 1 > i4; i4++)
                                        {
                                            tempName += " : " + tempSplits[i4];
                                        }
                                        if (sName[sName.Length - 3] == tempName)
                                        {
                                            rList.Add(i3);
                                        }
                                    }
                                }
                                else
                                    return rList;
                            }
                        }
                    }
                    else
                        return rList;
                }
            }

            return rList;
        }

        private int CreateBaseNodes(string[] labelSplits)
        {
            int i;
            for (i = 0; i < tvActiveWindows.Nodes.Count; i++)
            {
                if (tvActiveWindows.Nodes[i].Text.Split(new string[] { " : " },  StringSplitOptions.None)[0] == labelSplits[labelSplits.Length - 1])
                {
                    return  i;
                }
            }
            tvActiveWindows.Nodes.Add(labelSplits[labelSplits.Length - 1] + " :  0:00:00");
            return i;
        }

        private int CreateSubNodes(string[] labelSplits, int hiIndex)
        {
            int i;
            for (i = 0; i < tvActiveWindows.Nodes[hiIndex].Nodes.Count; i++)
            {
                if (tvActiveWindows.Nodes[hiIndex].Nodes[i].Text.Split(new string[] { " : " }, StringSplitOptions.None)[0] == labelSplits[labelSplits.Length - 2])
                {
                    return i;
                }
            }
            tvActiveWindows.Nodes[hiIndex].Nodes.Add(labelSplits[labelSplits.Length - 2] + " : 0:00:00");
            return i;
        }

        private int CreateSubSubNodes(string[] labelSplits, int hiIndex, int hiSubIndex)
        {
            int i;
            for (i = 0; i < tvActiveWindows.Nodes[hiIndex].Nodes[hiSubIndex].Nodes.Count; i++)
            {
                if (tvActiveWindows.Nodes[hiIndex].Nodes[hiSubIndex].Nodes[i].Text.Split(new string[] { " : " }, StringSplitOptions.None)[0] == labelSplits[labelSplits.Length - 3])
                {
                    return i;
                }
            }
            tvActiveWindows.Nodes[hiIndex].Nodes[hiSubIndex].Nodes.Add(labelSplits[labelSplits.Length - 3] + " : 0:00:00");
            return i;
        }

        private void FillNodes(string[] labelSplits, int hiIndex)
        {
            int holdvalue = 0;
            int holdvalue2 = 0;
            for (int i = labelSplits.Length - 2; i > -1; i--)
            {
                if ((labelSplits.Length - 2) == i)
                {
                    holdvalue = CreateSubNodes(labelSplits, hiIndex);
                }
                else if ((labelSplits.Length - 3) == i)
                {
                    holdvalue2 = CreateSubSubNodes(labelSplits, hiIndex, holdvalue);
                }
                else
                {
                    string temp = tvActiveWindows.Nodes[hiIndex].Nodes[holdvalue].Nodes[holdvalue2].Text;
                    tvActiveWindows.Nodes[hiIndex].Nodes[holdvalue].Nodes[holdvalue2].Text = (labelSplits[i] + " - " + temp);
                }
            }
        }

        private string FormatTime(int t)
        {
            t = t / 1000;
            int seconds = (t % 60);
            string SS = seconds.ToString();
            if (seconds < 10)
            {
                SS = '0' + SS;
            }
            int minutes = (((t - seconds) / 60) % 60);
            string SM = minutes.ToString();
            if (minutes < 10)
            {
                SM = '0' + SM;
            }
            int hours = ((t - minutes * 60 - seconds) / 3600);
            string SH = hours.ToString();
            return SH + ":" + SM + ":" + SS;
        }

        private int DeformatTime(string s)
        {
            string[] splits = s.Split(new string[] { ":" }, StringSplitOptions.None);
            int hours = int.Parse(splits[0]);
            int minutes = int.Parse(splits[1]);
            int seconds = int.Parse(splits[2]);
            int miliTime = 0;
            miliTime += (seconds * 1000);
            miliTime += (minutes * 60000);
            miliTime += (hours * 3600000);
            return miliTime;
        }

        private bool GetIsAFK()
        {
            int idleTime = 0;
            LASTINPUTINFO lastInputInfo = new LASTINPUTINFO();
            lastInputInfo.cbSize = Marshal.SizeOf(lastInputInfo);
            lastInputInfo.dwTime = 0;

            int envTicks = Environment.TickCount;

            if (GetLastInputInfo(out lastInputInfo))
            {
                int lastInputTick = lastInputInfo.dwTime;
                idleTime = envTicks - lastInputTick;
            }

            if (idleTime >= AFKTimeAmout)
            {
                return true;
            }
            else
                return false;
        }

        private void timer_Tick(object sender, EventArgs e)
        {
            //Finds which window is active
            GetActiveWindow();
            //label1.Text = "";
            //Adds to active window's time
            //for (int i = 0; i < activeWindowLabels.Count; i++)
            //{
            //    label1.Text += activeWindowLabels[i] + " - " + activeWindowTimes[i] + "\n";
            //}
            ShowActiveWindowsInTreeVeiw();
            //Checks if AFK
            isAFK = GetIsAFK();
            //Sets Total time
            TotalTimeWorking += timer.Interval;
            lblTime.Text = FormatTime(TotalTimeWorking);
        }

        private void saveToolStripMenuItem_Click(object sender, EventArgs e)
        {
            SaveFile();
        }

        private void openToolStripMenuItem_Click(object sender, EventArgs e)
        {
            OpenNewFile();
        }

        private void newToolStripMenuItem_Click(object sender, EventArgs e)
        {
            this.Text = "Time Keeper 2";
            lblProjectName.Text = "";
            SaveFileText = "";
            fileLocation = "";
            NewFile();
            lblStartingTime.Text = "";
            lblEndingTime.Text = "";
            btnStart.Text = "Start";
            timer.Stop();
            openToolStripMenuItem.Enabled = true;
            btnReview.Enabled = true;
        }

        private void btnReview_Click(object sender, EventArgs e)
        {
            timer.Stop();
            InitializeVaribles();
            lblStartingTime.Text = "";
            lblEndingTime.Text = "";
            try
            {
                for (int i = 0; dgvEntries.Rows.Count > i; i++)
                {
                    if (dgvEntries[0, i].Selected)
                    {
                        string[] saveList = SaveEntriesText[i].Split(new string[] { "\r\n" }, StringSplitOptions.None);

                        //Starting Time
                        if (lblStartingTime.Text == "")
                        {
                            lblStartingTime.Text = saveList[0];
                        }
                        //Ending Time
                        lblEndingTime.Text = saveList[1];

                        int activeWindowStartingIndex = 0;
                        //Nots
                        tbNotes.Text += saveList[0] + "\r\n\r\n";
                        for (int ni = 2; saveList.Length > ni; ni++)
                        {
                            if (saveList[ni] == "-SN%^1-")
                            {
                                activeWindowStartingIndex = ni + 1;
                                break;
                            }
                            tbNotes.Text += saveList[ni] + "\r\n";
                        }
                        tbNotes.Text += "\r\n";

                        //Active Windows
                        for (int ni = activeWindowStartingIndex; saveList.Length - 2 > ni; ni += 2)
                        {
                            bool ismatching = false;
                            int eleTime = int.Parse(saveList[ni + 1]);
                            TotalTimeWorking += eleTime;
                            for (int mi = 0; activeWindowLabels.Count > mi; mi++)
                            {
                                if (activeWindowLabels[mi] == saveList[ni])
                                {
                                    activeWindowTimes[mi] += eleTime;
                                    ismatching = true;
                                    break;
                                }
                            }
                            if (!ismatching)
                            {
                                activeWindowLabels.Add(saveList[ni]);
                                activeWindowTimes.Add(eleTime);
                            }
                        }
                    }
                }
            }
            catch (Exception)
            {
                MessageBox.Show("File is corrupted.", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error, MessageBoxDefaultButton.Button1);
                newToolStripMenuItem_Click(null, EventArgs.Empty);
                return;
            }

            lblTime.Text = FormatTime(TotalTimeWorking);
            ShowActiveWindowsInTreeVeiw();
        }

        private void btnStart_Click(object sender, EventArgs e)
        {
            if (lblStartingTime.Text == "")
            {
                timer.Start();
                this.Text = "Time Keeper 2";
                lblProjectName.Text = "";
                SaveFileText = "";
                fileLocation = "";
                lblNotes.Text = "";
                NewFile();
                btnStart.Text = "Stop";
                openToolStripMenuItem.Enabled = false;
                btnReview.Enabled = false;
                lblStartingTime.Text = DateTime.Now.ToString();
            }
            else if (btnStart.Text == "Stop")
            {
                SaveFile();
                if (fileLocation == "")
                {
                    timer.Start();
                    return;
                }
                OpenFile();
                btnStart.Text = "Start";
                openToolStripMenuItem.Enabled = true;
                btnReview.Enabled = true;
                lblStartingTime.Text = "-";
            }
            else if (btnStart.Text == "Start")
            {
                OpenFile();
                timer.Start(); 
                btnStart.Text = "Stop";
                openToolStripMenuItem.Enabled = false;
                btnReview.Enabled = false;
                lblStartingTime.Text = DateTime.Now.ToString();
            }
        }
    }
}
