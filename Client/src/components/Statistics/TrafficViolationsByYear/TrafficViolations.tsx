import React, { useState } from 'react';
import { Container, Form, Title, ErrorMessage, Loader, Table, TableHeader, TableRow, TableData } from './TrafficViolations.styled';
import { getTrafficViolationsByYear } from '../../../services/StatisticsService';
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';
import Button from '../../Shared/Button/Button';
import Input from '../../Shared/Input/Input';

interface TrafficViolationsReport {
    [key: string]: number;
}

const TrafficViolations: React.FC = () => {
    const [year, setYear] = useState('');
    const [violationsReport, setViolationsReport] = useState<TrafficViolationsReport | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [searchYear, setSearchYear] = useState<string | null>(null);

    const handleSearch = async (event: React.FormEvent) => {
        event.preventDefault();
        if (!/^\d{4}$/.test(year)) {
            setError('Please enter a valid year');
            setViolationsReport(null);
            setSearchYear(null);
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const response = await getTrafficViolationsByYear(year);
            setViolationsReport(response);
            setSearchYear(year);
            if (Object.values(response).every(count => count === 0)) {
                setError(`No traffic violations found for the year ${year}`);
            } else {
                setError(null);
            }
        } catch (error) {
            setError('Error fetching traffic violations');
            setViolationsReport(null);
        }

        setLoading(false);
    };

    const generatePDF = () => {
        const input = document.getElementById('report-content');
        if (input) {
            html2canvas(input).then((canvas) => {
                const imgData = canvas.toDataURL('image/png');
                const pdf = new jsPDF();
                const imgWidth = 190;
                const pageHeight = 295;
                const imgHeight = (canvas.height * imgWidth) / canvas.width;
                let heightLeft = imgHeight;

                let position = 10;

                pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
                heightLeft -= pageHeight;

                while (heightLeft >= 0) {
                    position = heightLeft - imgHeight;
                    pdf.addPage();
                    pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
                    heightLeft -= pageHeight;
                }

                pdf.save(`Traffic_Violations_Report_${searchYear}.pdf`);
            });
        }
    };

    // Proveri da li postoje podaci u violationsReport
    const isDownloadDisabled = !violationsReport || Object.values(violationsReport).every(count => count === 0);

    return (
        <Container>
            <Title>Traffic Violations</Title>
            <Form onSubmit={handleSearch}>
                <Input
                    type="text"
                    id="year"
                    label="Year"
                    attrName="year"
                    handleChange={(e: React.ChangeEvent<HTMLInputElement>) => setYear(e.target.value)}
                    data={year}
                />
                <Button label="Fetch Violations" buttonType="submit" />
            </Form>
            {loading && <Loader>Loading...</Loader>}
            {error && searchYear && <ErrorMessage>{error}</ErrorMessage>}
            {violationsReport && searchYear && Object.values(violationsReport).some(count => count > 0) && (
                <>
                    <div id="report-content">
                        <Table>
                            <thead>
                                <TableRow>
                                    <TableHeader>Reason</TableHeader>
                                    <TableHeader>Number of Violations</TableHeader>
                                </TableRow>
                            </thead>
                            <tbody>
                                {Object.entries(violationsReport).map(([reason, count]) => (
                                    <TableRow key={reason}>
                                        <TableData>{reason}</TableData>
                                        <TableData>{count}</TableData>
                                    </TableRow>
                                ))}
                            </tbody>
                        </Table>
                    </div>
                    <Button label="Download PDF" buttonType="button" onClick={generatePDF} disabled={isDownloadDisabled} />
                </>
            )}
        </Container>
    );
};

export default TrafficViolations;
