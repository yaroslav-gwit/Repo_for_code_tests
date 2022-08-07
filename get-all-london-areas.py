import requests
from bs4 import BeautifulSoup
import re

def get_underground() -> list:
    page = requests.get("https://en.wikipedia.org/wiki/List_of_London_Underground_stations")
    soup = BeautifulSoup(page.content, "html.parser")

    table = soup.find("table", class_="sortable")

    station_list = []
    for row in table.tbody.find_all('tr'):
        columns = row.find_all('th')
        if(columns != []):
            station = columns[0].text.strip()
            if station != "Station":
                station_list.append(station)
    return station_list


def get_overground() -> list:
    page = requests.get("https://en.wikipedia.org/wiki/Category:Railway_stations_served_by_London_Overground")
    soup = BeautifulSoup(page.content, "html.parser")
    soup_children = list(soup.children)
    body_text = soup_children[2].get_text()

    body_text_list_initial = body_text.split("\n")
    body_text_list_final = []
    for i in body_text_list_initial:
        if re.match(".*station$", i):
            i = i.replace(" railway station", "")
            i = i.replace(" station", "")
            body_text_list_final.append(i)
    
    station_list = body_text_list_final
    return station_list


def get_rails() -> list:
    page = requests.get("https://en.wikipedia.org/wiki/List_of_London_railway_stations")
    soup = BeautifulSoup(page.content, "html.parser")
    table = soup.find('table', class_='wikitable sortable')

    station_list = []
    for row in table.tbody.find_all('tr'):
        columns = row.find_all('td')
        if(columns != []):
            station = columns[0].text.strip()
            station = re.sub("\[.*\]", "", station)
            station_list.append(station)
    return station_list


def filter_out(input_string:str) -> str:
    i_split = re.split("([a-z]London)", input_string, 2)
    if len(i_split) > 1:
        partial = i_split[1].replace("London", "")
        i_split.pop(1)
        i_split.insert(1, partial)
        result = i_split[0] + i_split[1]
    else:
        result = i_split[0]
    return result


def final_list() -> list:
    underground_stations = get_underground()
    overground_stations = get_overground()
    rail_stations = get_rails()

    final_list = []
    for i in underground_stations:
        i = filter_out(i)
        final_list.append(i)
    for i in overground_stations:
        i = filter_out(i)
        final_list.append(i)
    for i in rail_stations:
        i = filter_out(i)
        final_list.append(i)

    _set = set(final_list)
    final_list = list(_set)
    final_list.sort()
    return final_list


if __name__ == "__main__":
    _list = final_list()
    print(_list)
    print(len(_list))
